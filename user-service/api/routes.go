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
	mux.HandleFunc("/users/enable/user", handlers.EnableUser)
	mux.HandleFunc("/users/disable/user", handlers.DisableUser)
	mux.HandleFunc("/users/retrieve/user", handlers.RetrieveUserUsername)
	mux.HandleFunc("/users/retrieve/type", handlers.RetrieveUserType)

	mux.HandleFunc("/users/audit/logs", handlers.GetUserAuditActions)
	mux.HandleFunc("/users/list/users", handlers.RetrieveUsers)
}
