package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID            uuid.UUID   `json:"id"`
	Resource      uuid.UUID   `json:"resource"`    // The resource the permission applies to (e.g., "warehouses")
	ResourceID    string      `json:"resource_id"` // The UUID of the resource
	ResourceName    string      `json:"resource_name"` // The na,e of the resource
	Operations    []uuid.UUID `json:"operations"`
	OperationName string      `json:"operation_name"`
	Effect        string      `json:"effect"`     // The effect of the permission (e.g., "allow", "deny")
	CreatedAt     time.Time   `json:"created_at"` // The timestamp of when the permission was created
	UpdatedAt     time.Time   `json:"updated_at"` // The timestamp of the last update to the permission
}
