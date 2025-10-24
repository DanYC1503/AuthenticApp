package api

import (
	"main/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {

	// User routes
	mux.HandleFunc("/audit/log", handlers.AuditAction)
	

}
