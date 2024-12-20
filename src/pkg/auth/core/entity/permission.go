package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID                 uuid.UUID `json:"id"`
	Resource           string    `json:"resource"`            // The resource the permission applies to (e.g., "warehouses")
	ResourceIdentifier string    `json:"resource_identifier"` // The identifier for the specific resource (e.g., warehouse ID)
	Operation          string    `json:"operation"`           // The operation allowed (e.g., "read", "write", "delete")
	Effect             string    `json:"effect"`              // The effect of the permission (e.g., "allow", "deny")
	CreatedAt          time.Time `json:"created_at"`          // The timestamp of when the permission was created
	UpdatedAt          time.Time `json:"updated_at"`          // The timestamp of the last update to the permission
}
