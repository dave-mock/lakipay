package entity

import (
	"time"

	"github.com/google/uuid"
)

type LocationType string

const (
	GEO      LocationType = "GEO"
	RELATIVE LocationType = "RELATIVE"
	LOCAL    LocationType = "LOCAL"
)

type Address struct {
	Id        uuid.UUID
	Title     string
	Phones    []string
	Emails    []string
	Websites  []string
	Locations []Location
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Location struct {
	Type    LocationType
	Details interface{}
}

/*
A geo location is a representation of location in the form
of a standard Latitude and Longutude notation
*/

type GeoLocation struct {
	Lat float64
	Lng float64
}

/*
A relative location is a representation of a location data
using a known place around where your business is setup
*/
type RelativeLocation struct {
	Direction string
}

/*
Local location represents a location data collected in
the form of a formal representation in the provided country
*/
type LocalLocation struct {
	Data map[string]string
}
