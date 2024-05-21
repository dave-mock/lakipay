package entity

import (
	"image/color"

	"github.com/google/uuid"
)

type Tag struct {
	Id    uuid.UUID
	Name  string
	Color color.RGBA64
}
