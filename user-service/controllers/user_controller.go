package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/config"
	"main/middleware"
	"main/models"
	repository "main/repository/user-processing"

	"net/http"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	const maxRetries = 3
	db := config.ConnectDB()
	fmt.Println("Connected to Database")
	defer db.Close()
	//  Decode update payload
	var user models.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Pass a function (closure) that receives *tx and calls your repository function
	middleware.TransactionRetry(db, maxRetries, w, func(tx *sql.Tx) error {
		return repository.UpdateUser(w, user, tx)
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User updated successfully")
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	const maxRetries = 3
	db := config.ConnectDB()
	fmt.Println("Connected to Database")
	defer db.Close()

	var userInfo models.UserRequestInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Pass a function (closure) that receives *tx and calls your repository function
	middleware.TransactionRetry(db, maxRetries, w, func(tx *sql.Tx) error {
		return repository.RetrieveUserLoginInfo(w, userInfo, tx)
	})

	// Success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Info returned successfully"))
}
