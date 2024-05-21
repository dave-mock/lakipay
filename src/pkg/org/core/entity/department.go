package entity

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	Id          uuid.UUID
	Name        string
	Description string
	Logo        string
	RegDate     time.Time
	Addresses   []Address
	Categories  []Category
	Services    []Service
	Offers      []Offer
	Details     interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Fashion Retailer
//
// Software Development
type Service struct {
	Id               uuid.UUID
	Name             string
	Description      string
	Categories       []Category
	CountryWhitelist []string
	CountryBlacklist []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Nike AirForce 1
// Hp Envy 17
// Mobile App Development | Website Development
type Offer struct {
	Id uuid.UUID
}

/*

: Title : Nike AirForce 1
: Type  : Shoe
	- Men's
	- Trainer | Sneaker | Casual
	- Size
	- Color
	- Fabric

	* Option

	*

	Node -> Leaf

*/
