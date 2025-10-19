package api

import (
	"main/handlers"

	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {

	// User routes
	mux.HandleFunc("/user/update", handlers.UpdateUser)
	mux.HandleFunc("/user/info", handlers.GetUserInfo)

}
