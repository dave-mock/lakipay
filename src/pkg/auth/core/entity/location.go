package entity

import "github.com/google/uuid"

// Location
type Location struct {
	Id       uuid.UUID
	Relative *RelativeLocation
	Geo      *GeoLocation
}

type RelativeLocation struct {
	Directions string
}

type GeoLocation struct {
	Lat float64
	Lng float64
}
