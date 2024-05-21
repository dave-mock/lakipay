package psql

import (
	"net/http"
)

type PsqlRepo struct{
	// Logger
	// log slog.
	sm *http.ServeMux
	
}