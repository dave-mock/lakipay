package entity

import "github.com/google/uuid"

type VerificationStatus struct {
	Id       uuid.UUID
	Verified bool
	Status   string
	Message  string
}
