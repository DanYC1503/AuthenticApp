package api

import (
	"main/controllers"
	"main/handlers"

	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {

	// User routes
	mux.HandleFunc("/auth/register", handlers.RegisterUser)
	mux.HandleFunc("/auth/login", handlers.LoginUser)

	// OAuth Paths
	mux.HandleFunc("/auth/google/login", handlers.GoogleLogin)
	mux.HandleFunc("/auth/google/callback", handlers.GoogleCallback)

	//mux.HandleFunc("GET /auth/refresh", handlers.RefreshToken)
	mux.HandleFunc("/auth/validateToken", handlers.TokenVerification)
	mux.HandleFunc("/auth/validateUpToken", handlers.UpTokenVerification)
	mux.HandleFunc("/auth/validateDelToken", handlers.DelTokenVerification)
	mux.HandleFunc("/auth/validatePasswordToken", handlers.ResetPasswordVerification)

	mux.HandleFunc("/auth/deleteToken", controllers.RequireValidToken(handlers.GetDeleteToken))
	mux.HandleFunc("/auth/updateUserToken", controllers.RequireValidToken(handlers.GetUpdateToken))
	mux.HandleFunc("/auth/logout", controllers.RequireValidToken(handlers.LogoutSession))
	mux.HandleFunc("/auth/password/reset", handlers.ResetPassword)
	mux.HandleFunc("/auth/password/token", handlers.GetPasswordToken)

}
