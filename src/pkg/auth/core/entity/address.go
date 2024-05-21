package entity

import "github.com/google/uuid"

type Address struct {
	Id       uuid.UUID
	Title    string
	Primary  bool
	Phones   []Phone
	Emails   []string
	Location Location
}
