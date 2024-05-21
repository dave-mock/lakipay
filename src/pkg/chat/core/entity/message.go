package entity

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	PHOTO    MediaType = "PHOTO"
	GEO      MediaType = "GEO"
	CONTACT  MediaType = "CONTACT"
	DOCUMENT MediaType = "DOCUMENT"
	WEB_PAGE MediaType = "WEB_PAGE"
	INVOICE  MediaType = "INVOICE"
)

type DocumentAttribute string

const (
	IMAGE DocumentAttribute = "IMAGE"
	AUDIO DocumentAttribute = "AUDIO"
	VIDEO DocumentAttribute = "VIDEO"
	FILE  DocumentAttribute = "FILE"
)

type PeerType string

const (
	USER    PeerType = "USER"
	GROUP   PeerType = "GROUP"
	CHANNEL PeerType = "CHANNEL"
	SYSTEM  PeerType = "SYSTEM"
)

type Peer struct {
	Id   uuid.UUID
	Type PeerType
}

type Message struct {
	Id        uuid.UUID
	From      Peer
	To        Peer
	Message   string
	Media     *MessageMedia
	CreatedAt time.Time
	UpdatedAt time.Time
	GroupId   uuid.UUID
}

type MessageMedia struct {
	Id         uuid.UUID
	Type       MediaType
	AccessHash string
	DocId      string
	Thumbnail  []int8
	Details    interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MediaDocument struct {
	MIMEType  string
	Attribute DocumentAttribute
	Size      float64
}

type MediaPhoto struct {
	Width  float64
	Height float64
	Size   float64
}
