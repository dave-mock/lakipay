package entity

import (
	"time"

	"github.com/google/uuid"
)

// Object data and Metadata,
// The metadata is a set of name-value pairs that describe the object.
// These pairs include some default metadata, such as the date last modified, and standard HTTP metadata, such as Content-Type.
// You can also specify custom metadata at the time that the object is stored.
type Object struct {
	Id        uuid.UUID
	Key       string
	VersionId uuid.UUID
	Value     []byte
	Metadata  Metadata
}

type Metadata struct {
	Size float64
	Date time.Time
	// A general header field used to specify caching policies.
	CacheControl interface{}
	// Object presentational information.
	ContentDisposition interface{}
	// The object size in bytes.
	ContentLength interface{}
	// The object type.
	ContentType interface{}
	// The object creation date or the last modified date, whichever is the latest.
	// For multipart uploads, the object creation date is the date of initiation of the multipart upload.
	LastModified time.Time
	// MD5 digest of the data.
	ETag interface{}
}
