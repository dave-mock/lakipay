package etrade

import (
	"auth/src/pkg/org/core/entity"
	"auth/src/pkg/org/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
)

var legalConditions map[string]string = map[string]string{
	"1":  "Private",
	"2":  "Private Limited Company",
	"3":  "Share Company",
	"4":  "Commercial Representative",
	"5":  "Public EnterPrise",
	"6":  "Paretnership",
	"7":  "Cooperatives Association",
	"8":  "Trade Sectoral Association",
	"9":  "Non Public EnterPrise",
	"10": "NGO",
	"11": "Branch of A foreign Chamber of Commerce",
	"12": "Holding Company",
	"13": "Franchising",
	"14": "Border Trade",
	"15": "International Bid Winners Foreign Companies",
	"16": "One Man Private Limited Company",
}

type Etrade struct {
	log     *log.Logger
	Referer string
}

func New(log *log.Logger) usecase.TINChecker {
	return Etrade{log: log, Referer: "https://etrade.gov.et/business-license-checker"}
}

func (etrade Etrade) CheckTIN(tin string, uc usecase.Usecase) (*entity.Organization, error) {

	etrade.log.Println("CHECKING TIN")

	type Bus struct {
		LicenceNumber              string  `json:"LicenceNumber"`
		Status                     int64   `json:"Status"`
		Capital                    float64 `json:"Capital"`
		RenewedTo                  string  `json:"RenewedTo"`
		BusinessLicensingGroupMain []struct {
			SubGroup     int64 `json:"SubGroup"`
			BusinessMain struct {
				DateRegistered string `json:"DateRegistered"`
				TradesName     string `json:"TradesName"`
				RenewedTo      string `json:"RenewedTo"`
			} `json:"BusinessMain"`
		} `json:"BusinessLicensingGroupMain"`
	}

	type Org struct {
		Tin           string  `json:"Tin"`
		LegalCondtion string  `json:"LegalCondtion"`
		RegNo         string  `json:"RegNo"`
		RegDate       string  `json:"RegDate"`
		BusinessName  string  `json:"BusinessName"`
		PaidUpCapital float64 `json:"PaidUpCapital"`
		Businesses    []struct {
			OwnerTIN      string `json:"OwnerTIN"`
			LicenceNumber string `json:"LicenceNumber"`
		} `json:"Businesses"`
	}

	// Request for Org Info
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://etrade.gov.et/api/Registration/GetRegistrationInfoByTin/%s/en", tin), nil)
	if err != nil {
		log.Println("CHECKING TIN ERROR REQ")
		return nil, err
	}

	req.Header.Add("Referer", etrade.Referer)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("failed to check TIN")
	}

	var org Org

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&org)
	if err != nil {
		return nil, err
	}

	etrade.log.Println("ORG DONE")
	etrade.log.Println(org.BusinessName)
	etrade.log.Println(org.PaidUpCapital)
	etrade.log.Println(org.RegDate)
	etrade.log.Println(org.LegalCondtion)

	regDate, _ := time.Parse("1/2/2006", org.RegDate)
	category, _ := uc.GetCategoryByName("Business")
	legalCondition, _ := uc.GetLegalConditionByName(legalConditions[org.LegalCondtion])

	etrade.log.Println(regDate)
	etrade.log.Println("category")
	etrade.log.Println(category)
	etrade.log.Println(legalCondition)

	var parsedOrg entity.Organization = entity.Organization{
		Id:             uuid.Nil,
		Name:           org.BusinessName,
		Capital:        org.PaidUpCapital,
		RegDate:        regDate,
		Country:        "ET",
		Category:       category,
		LegalCondition: legalCondition,
		Details: entity.EthBusOrg{
			TIN:   org.Tin,
			RegNo: org.RegNo,
		},
	}

	//////////////////////////////// ORG DONE ////////////////////////////////////

	// var buses []Bus

	type Cat struct {
		Id         int64
		Name       string
		RegDate    string
		ValidUntil string
		LicenseNo  string
	}

	var _buses map[string][]Cat = make(map[string][]Cat)

	log.Println("BUS START")
	log.Println(len(org.Businesses))

	for i := 0; i < len(org.Businesses); i++ {
		log.Println("BUS START - LOOPING")
		// Request Businesses
		req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://etrade.gov.et/api/BusinessMain/GetBusinessByLicenseNo?LicenseNo=%s&Tin=%s&Lang=en", org.Businesses[i].LicenceNumber, org.Tin), nil)
		if err != nil {
			log.Println("CHECKING business ERROR REQ")
			return nil, err
		}

		log.Println("BUS START - 1")

		req.Header.Add("Referer", etrade.Referer)

		client = http.Client{}
		res, err = client.Do(req)
		if err != nil {
			log.Println("BUS START - 2")
			return nil, err
		}

		log.Println("BUS START - 3")

		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Println("BUS START - 4")
			return nil, errors.New("failed to check business")
		}

		var bus Bus

		log.Println("BUS START - 5")

		decoder = json.NewDecoder(res.Body)
		err = decoder.Decode(&bus)
		if err != nil {
			log.Println(err.Error())
			log.Println("BUS START - 6")
			return nil, err
		}

		log.Println("BUS START - 7")

		if len(bus.BusinessLicensingGroupMain) > 0 {
			log.Println(bus.BusinessLicensingGroupMain[0].BusinessMain.TradesName)

			if reflect.DeepEqual(_buses[bus.BusinessLicensingGroupMain[0].BusinessMain.TradesName], Cat{}) {
				_buses[bus.BusinessLicensingGroupMain[0].BusinessMain.TradesName] = []Cat{
					{
						Id:         bus.BusinessLicensingGroupMain[0].SubGroup,
						Name:       "PENDING PARSE",
						RegDate:    bus.BusinessLicensingGroupMain[0].BusinessMain.DateRegistered,
						ValidUntil: bus.BusinessLicensingGroupMain[0].BusinessMain.RenewedTo,
						LicenseNo:  bus.LicenceNumber,
					},
				}
			} else {
				_buses[bus.BusinessLicensingGroupMain[0].BusinessMain.TradesName] = append(_buses[bus.BusinessLicensingGroupMain[0].BusinessMain.TradesName], Cat{
					Id:         bus.BusinessLicensingGroupMain[0].SubGroup,
					Name:       "PENDING PARSE",
					RegDate:    bus.BusinessLicensingGroupMain[0].BusinessMain.DateRegistered,
					ValidUntil: bus.BusinessLicensingGroupMain[0].BusinessMain.RenewedTo,
					LicenseNo:  bus.LicenceNumber,
				})

			}
		}

	}

	log.Println("BUS DONE LOOPING")
	log.Println(len(_buses))

	var departments []entity.Department

	for k, v := range _buses {

		var cats []entity.Category = make([]entity.Category, 0)

		for _, c := range v {
			cats = append(cats, entity.Category{
				// Id:   c.Id,
				Name: c.Name,
			})
		}

		// regDate, _ := time.Parse("2006-01-02T15:04:05", buses[i].BusinessLicensingGroupMain[0].BusinessMain.DateRegistered)
		// valid, _ := time.Parse("2006-01-02T15:04:05", buses[i].BusinessLicensingGroupMain[0].BusinessMain.RenewedTo)

		departments = append(departments, entity.Department{
			Name:       k,
			Categories: cats,
			Services:   make([]entity.Service, 0),
			Offers:     make([]entity.Offer, 0),
		})
	}

	parsedOrg.Departments = departments

	return &parsedOrg, nil
}

// TEST TINs
// 0000030603
// 0014707151
// 0073327705

// Refresh Token
// Report
// Trace Back

// Multiple Associates

//
