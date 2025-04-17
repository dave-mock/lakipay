package rest

import (
	"auth/src/pkg/account/core/entity"
	"net/http"
	"time"
)

var timeSnap, _ = time.Parse("2006-01-02 15:04:05", "2025-03-20 11:34:12")

var gatewayList []map[string]interface{} = []map[string]interface{}{
	// [TELEBIRR]
	encodeGateway(entity.Gateway{
		Id:      "D1D47D97-101A-4DCC-965A-6E7C236BC2EE",
		Key:     "TELEBIRR",
		Name:    "Telebirr",
		Acronym: "",
		Icon:    "https://play-lh.googleusercontent.com/Mtnybz6w7FMdzdQUbc7PWN3_0iLw3t9lUkwjmAa_usFCZ60zS0Xs8o00BW31JDCkAiQk",

		Type:       entity.WALLET,
		CanProcess: true,
		CanSettle:  true,
		CreatedAt:  timeSnap,
		UpdatedAt:  timeSnap,
	}),

	// [CBE BIRR]
	encodeGateway(entity.Gateway{
		Id:      "FD5D56EF-061D-432A-A93C-C57BA3061914",
		Key:     "CBE_BIRR",
		Name:    "CBE Birr",
		Acronym: "",
		Icon:    "https://play-lh.googleusercontent.com/rcSKabjkP2GfX1_I_VXBfhQIPdn_HPXj5kbkDoL4cu5lpvcqPsGmCqfqxaRrSI9h5_A",

		Type:       entity.WALLET,
		CanProcess: true,
		CanSettle:  true,
		CreatedAt:  timeSnap,
		UpdatedAt:  timeSnap,
	}),

	// [CYBERSOURCE]
	encodeGateway(entity.Gateway{
		Id:      "0EC7A9C8-B40D-4DC8-9082-5A092697958A",
		Key:     "CYBERSOURCE",
		Name:    "VISA MASTERCARD",
		Acronym: "",
		Icon:    "https://getsby.com/wp-content/uploads/2023/01/Visa-Mastercard-1-1024x378.png",

		Type:       entity.CARD,
		CanProcess: true,
		CanSettle:  false,
		CreatedAt:  timeSnap,
		UpdatedAt:  timeSnap,
	}),

	/// [BANK]
	///
	/// [AWASH]
	encodeGateway(entity.Gateway{
		Id:      "B078C9AE-045B-4D3E-981F-17D46A1E8F75",
		Key:     "AWINETAA",
		Name:    "Awash International Bank SC",
		Acronym: "Awash",
		Icon:    "https://upload.wikimedia.org/wikipedia/commons/3/33/Awash_International_Bank.png",

		Type:       entity.BANK,
		CanProcess: true,
		CanSettle:  true,
		CreatedAt:  timeSnap,
		UpdatedAt:  timeSnap,
	}),
	///
	/// [BUNNA]
	encodeGateway(entity.Gateway{
		Id:      "1B2D3252-C25A-4952-BB24-30919FF23C94",
		Key:     "BUNAETAA",
		Name:    "Bunna International Bank SC",
		Acronym: "Bunna",
		Icon:    "https://z-p3-scontent.fadd1-1.fna.fbcdn.net/v/t39.30808-6/272417075_4734443969964854_7325903198322919211_n.jpg?_nc_cat=103&ccb=1-7&_nc_sid=6ee11a&_nc_ohc=Bh4lLYSzbAgQ7kNvgGM6iHo&_nc_oc=AdlnWv8iU5J5E-N4gM0VeNpQsGoDAN1L7yLc8JiLvXdqcxti2w4FQNciLOIH9FysaNg&_nc_zt=23&_nc_ht=z-p3-scontent.fadd1-1.fna&_nc_gid=-wTAXo6UP9h1fEn6JMygEg&oh=00_AYHzaOpQuXSqV4XfHO2bmD95rregJSgUbwhHjJTD1kH-Og&oe=67E1A7AB",

		Type:       entity.BANK,
		CanProcess: true,
		CanSettle:  true,
		CreatedAt:  timeSnap,
		UpdatedAt:  timeSnap,
	}),
}

func encodeGateway(v entity.Gateway) map[string]interface{} {
	return map[string]interface{}{
		"id":         v.Id,
		"key":        v.Key,
		"name":       v.Name,
		"short_name": v.Acronym,
		"icon":       v.Icon,

		"type":        v.Type,
		"can_process": v.CanProcess,
		"can_settle":  v.CanSettle,
		"created_at":  v.CreatedAt,
		"updated_at":  v.UpdatedAt,
	}
}

func (controller Controller) getGateways(w http.ResponseWriter, _ *http.Request) {

	SendJSONResponse(w, Response{
		Success: true,
		Data:    gatewayList,
	}, 200)
}
