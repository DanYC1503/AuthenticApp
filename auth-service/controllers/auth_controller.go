package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"main/config"
	"main/middleware"
	"main/middleware/encryption"
	"main/models"
	repository "main/repository/auth-processing"

	"net/http"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserCreateClient
	const maxRetries = 3
	var createDate = time.Now()
	db := config.ConnectDB()
	fmt.Println("Connected to Database")
	defer db.Close()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest)
		return
	}
	fmt.Println("Email: " + user.Email)
	fmt.Println("FullName: " + user.FullName)
	fmt.Println("Username: " + user.Username)
	if user.Email == "" || user.FullName == "" || user.Username == "" {
		http.Error(
			w,
			"Name, Email and Username are required",
			http.StatusBadRequest,
		)
		return

	}
	fmt.Printf("In controllers trying to insert user")
	for attempt := 1; attempt <= maxRetries; attempt++ {
		tx, err := db.Begin()
		if err != nil {
			log.Println("Begin failed:", err)
			continue
		}
		fmt.Println("Attempting to insert user")
		errIns := repository.InsertUser(tx, user, createDate)
		if errIns != nil {
			tx.Rollback()
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		if err != nil {
			tx.Rollback()
			if encryption.IsRetryable(err) {
				log.Printf("Retry %d/%d: %v\n", attempt, maxRetries, err)
				time.Sleep(time.Duration(attempt) * time.Second) // backoff
				continue
			}
			return
		}

		if err := tx.Commit(); err != nil {
			if encryption.IsRetryable(err) {
				log.Printf("Retry %d/%d: Commit failed: %v\n", attempt, maxRetries, err)
				continue
			}
			return
		}

		w.Header().Set("X-Username", user.Username)
		w.WriteHeader(http.StatusCreated)
		break
	}
}
func ParseLoginRequest(r *http.Request) (models.UserLogin, error) {
	var user models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return models.UserLogin{}, err
	}
	return user, nil
}
func LoginUser(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	// Parse login request
	user, err := ParseLoginRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Authenticate credentials
	ok, username, err := repository.AuthenticateUserCredentials(db, user.Username, user.Password)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate session token
	sessionToken, sessionExp, err := encryption.GenerateToken(username, "session")
	if err != nil {
		http.Error(w, "Could not generate session token", http.StatusInternalServerError)
		return
	}
	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  sessionExp,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // enable in production
		SameSite: http.SameSiteStrictMode,
	})

	// Respond with session details only

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Username", user.Username)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_token": sessionToken,
		"expires":       sessionExp.Unix(),
	})
}

func GetDeleteToken(w http.ResponseWriter, r *http.Request) {

	var user models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token := repository.RequestDeleteAuthToken(w, r, user)
	resp := map[string]string{"deleteAuthToken": token}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Username", user.Username)
	json.NewEncoder(w).Encode(resp)
}

func SessionTokenVerification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	isValid, username := repository.SessionTokenVerification(w, r)

	if !isValid {
		http.Error(w, "Invalid or missing session token", http.StatusUnauthorized)
		return
	}

	resp := map[string]string{
		"tokenStatus": "valid",
		"username":    username,
	}
	json.NewEncoder(w).Encode(resp)
}
func UpdateTokenVerification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	updateToken := r.Header.Get("X-Update-Auth")
	if updateToken == "" {
		http.Error(w, "Missing update token", http.StatusUnauthorized)
		return
	}

	claims, ok := encryption.ValidateToken(updateToken, "update")
	if !ok {
		http.Error(w, "Invalid or expired update token", http.StatusUnauthorized)
		return
	}

	status := fmt.Sprintf("Update token validated for user: %s", claims.Username)
	resp := map[string]string{"tokenStatus": status}
	json.NewEncoder(w).Encode(resp)
}

func DeleteTokenVerification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	deleteToken := r.Header.Get("X-Delete-Auth")
	if deleteToken == "" {
		http.Error(w, "Missing delete token", http.StatusUnauthorized)
		return
	}

	claims, ok := encryption.ValidateToken(deleteToken, "delete")
	if !ok {
		http.Error(w, "Invalid or expired delete token", http.StatusUnauthorized)
		return
	}

	status := fmt.Sprintf("Delete token validated for user: %s", claims.Username)
	resp := map[string]string{"tokenStatus": status}
	json.NewEncoder(w).Encode(resp)
}

func RequireValidToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isValid, _ := repository.SessionTokenVerification(w, r)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid or missing token"))
			return
		}
		// Token is valid, continue to the next handler
		next(w, r)
	}
}
func GetUpdateToken(w http.ResponseWriter, r *http.Request) {
	var user models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token := repository.RequestUpdateAuthToken(w, r, user)
	resp := map[string]string{"updateAuthToken": token}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Username", user.Username)
	json.NewEncoder(w).Encode(resp)
}

func LogoutSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Attempt to logout
	success := middleware.LogoutCurrentUser(w, r)

	status := "0"
	if success {
		status = "1"
	}

	resp := map[string]string{
		"session": "loggedout",
		"success": status,
	}

	json.NewEncoder(w).Encode(resp)
}
