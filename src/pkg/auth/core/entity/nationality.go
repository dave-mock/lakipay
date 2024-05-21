package entity

import "github.com/google/uuid"

type Nationality struct {
	Id          uuid.UUID
	Country     Country
	Description string
}
