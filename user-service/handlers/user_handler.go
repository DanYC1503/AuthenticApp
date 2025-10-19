package handlers

import (
	"main/controllers"

	"net/http"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	controllers.GetUserInfo(w, r)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	controllers.UpdateUser(w, r)
}
