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

	username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "Username query parameter is required", http.StatusBadRequest)
        return
    }

    // Create UserRequestInfo from query parameter
    userInfo := models.UserRequestInfo{
        Username: username,
    }

	// Pass a function (closure) that receives *tx and calls your repository function
	middleware.TransactionRetry(db, maxRetries, w, func(tx *sql.Tx) error {
		return repository.RetrieveUserLoginInfo(w, userInfo, tx)
	})

}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
		return repository.DeleteUser(w, userInfo, tx)
	})

}
func DisableUser(w http.ResponseWriter, r *http.Request) {
	const maxRetries = 3
	db := config.ConnectDB()
	fmt.Println("Connected to Database")
	defer db.Close()

	var userInfo models.AdminUserEnableDisable
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Pass a function (closure) that receives *tx and calls your repository function
	middleware.TransactionRetry(db, maxRetries, w, func(tx *sql.Tx) error {
		return repository.DisableUser(w, userInfo, tx)
	})

}
func EnableUser(w http.ResponseWriter, r *http.Request) {
	const maxRetries = 3
	db := config.ConnectDB()
	fmt.Println("Connected to Database")
	defer db.Close()

	var userInfo models.AdminUserEnableDisable
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Pass a function (closure) that receives *tx and calls your repository function
	middleware.TransactionRetry(db, maxRetries, w, func(tx *sql.Tx) error {
		return repository.EnableUser(w, userInfo, tx)
	})

}
func RetrieveUserUsername(w http.ResponseWriter, r *http.Request) {
	const maxRetries = 3
	db := config.ConnectDB()
	fmt.Println("Connected to Database")
	defer db.Close()

	var userInfo models.UserRequestInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Verify if user is admin before allowing access
	isAdmin, err := repository.VerifyUserType(db, userInfo.Username)
	if err != nil {
		http.Error(w, "Internal error while verifying admin", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "User not authorized to view other users", http.StatusForbidden)
		return
	}

	// Execute within transaction with retry
	middleware.TransactionRetry(db, maxRetries, w, func(tx *sql.Tx) error {
		return repository.RetrieveUserUsername(w, tx)
	})
}
func ReturnUserType(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()
	var userInfo models.UserRequestInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	userType, err := repository.ReturnUserType(tx, userInfo.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tx.Commit()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"type": userType,
	})
}
