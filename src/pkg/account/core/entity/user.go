package entity

import "github.com/google/uuid"

type User struct {
	Id uuid.UUID
}
type User2 struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
