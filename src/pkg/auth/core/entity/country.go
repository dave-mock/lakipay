package entity

import "github.com/google/uuid"

type Country struct {
	Id          uuid.UUID
	Name        string
	DefaultName string
	Flag        string
	Iso2        string
	Hidden      bool
	// [NB] - One country can only have a single phone prefix
	PhonePrefix PhonePrefix
}

type PhonePrefix struct {
	Prefix  string
	Pattern string
}
