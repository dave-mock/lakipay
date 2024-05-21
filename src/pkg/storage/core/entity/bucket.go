package entity

import "github.com/google/uuid"

type Bucket struct {
	Id                uuid.UUID
	Name              string
	DefaultEncryption bool
	Versioning        bool
	Objects           []Object
}
