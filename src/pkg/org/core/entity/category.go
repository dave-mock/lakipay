package entity

import (
	"time"

	"github.com/google/uuid"
)

// Retailers
//
// Technology
type Category struct {
	Id               uuid.UUID
	Name             string
	Description      string
	Icon             string
	Parents          []uuid.UUID
	CountryWhitelist []string
	CountryBlacklist []string
	// Add Specification Options
	Options   []Option
	Hidden    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Option

/*
Name:
Desc:
DataType:
Rep: // cm, gb, lt, ...
Values: []expected_values
Allow Custom Value: true
Validator:
Conditions (required, minLength, maxLength // React hook form)
*/

type DataType string

const (
	STRING  DataType = "string"
	INTEGER DataType = "integer"
	REAL    DataType = "real"
)

type Option struct {
	Id               uuid.UUID
	Name             string
	Description      string
	DataType         DataType
	RepresentedIn    string
	Values           []interface{}
	AllowCustomValue bool
	Validator        OptionValidator
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Validate option items using a predefined set of
// Keys and values
/*
e.g. required: {
		value: true,
		message: "Field is required"
	}
*/
type OptionValidator map[string]struct {
	Value   interface{}
	Message interface{}
}
