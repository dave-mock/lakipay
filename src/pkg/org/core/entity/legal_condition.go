package entity

import (
	"time"

	"github.com/google/uuid"
)

/*
Legal Condition
- Business forming types (PLC, SC, ...)
*/

type LegalCondition struct {
	Id               uuid.UUID
	Name             string
	Description      string
	CountryWhitelist []string
	CountryBlacklist []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

/*

   "Description": "Private",
   "Description": "Private Limited Company",
   "Description": "share Company",
   "Description": "Commercial Representative",
   "Description": "Public EnterPrise",
   "Description": "Paretnership",
   "Description": "Cooperatives Association",
   "Description": "Trade Sectoral Association",
   "Description": "Non Public EnterPrise",
   "Description": "NGO",
   "Description": "Branch of A foreign Chamber of Commerce",
   "Description": "Holding Company",
   "Description": "Franchising",
   "Description": "Border Trade",
   "Description": "International Bid Winners Foreign Companies",
   "Description": "One Man Private Limited Company",

*/

var LegalConditions []LegalCondition = []LegalCondition{
	// Private
	{
		Id:               uuid.New(),
		Name:             "Private",
		Description:      "Private owned organization",
		CountryWhitelist: []string{"ET"},
		CountryBlacklist: []string{},
		CreatedAt:        time.Now(),
	},
	// Private Limited Company
	{
		Id:               uuid.New(),
		Name:             "Private Limited Company",
		Description:      "PLC organization",
		CountryWhitelist: []string{"ET"},
		CountryBlacklist: []string{},
		CreatedAt:        time.Now(),
	},
	// Share Company
	{
		Id:               uuid.New(),
		Name:             "Share Company",
		Description:      "Share organization",
		CountryWhitelist: []string{"ET"},
		CountryBlacklist: []string{},
		CreatedAt:        time.Now(),
	},
	// Commercial Representative
	{
		Id:               uuid.New(),
		Name:             "Commercial Representative",
		Description:      "Commercial Representative organization",
		CountryWhitelist: []string{"ET"},
		CountryBlacklist: []string{},
		CreatedAt:        time.Now(),
	},
}
