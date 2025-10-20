package api

import (
	"main/handlers"

	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {

	// User routes
	mux.HandleFunc("/users/update", handlers.UpdateUser)
	mux.HandleFunc("/users/info", handlers.GetUserInfo)
	mux.HandleFunc("/users/delete", handlers.DeleteUser)

}
