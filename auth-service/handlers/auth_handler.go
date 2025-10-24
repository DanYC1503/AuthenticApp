package handlers

import (
	"context"
	"fmt"
	"log"
	"main/controllers"
	"main/middleware"
	"net/http"

	"github.com/markbates/goth/gothic"
)

// ---------------------Basic auth funtions
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	fmt.Printf("Register User reached going to controllers")
	controllers.CreateUser(w, r)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.LoginUser(w, r)
}

// ----------------------TOKEN RETRIEVAL
func GetDeleteToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.GetDeleteToken(w, r)
}
func GetUpdateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.GetUpdateToken(w, r)
}
func GetPasswordToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.GetPasswordToken(w, r)
}

// -------------------TOKEN VERIFICATION
func AutoVerifyRecoveryToken(w http.ResponseWriter, r *http.Request) {

	// Reuse your existing verification logic
	controllers.VerifyPasswordRecoveryToken(w, r)
}

func TokenVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	controllers.SessionTokenVerification(w, r)
}
func UpTokenVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.UpdateTokenVerification(w, r)
}
func DelTokenVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	deleteToken := r.Header.Get("X-Delete-Auth")
	log.Printf("Auth-controllers received delete token: %s", deleteToken)

	controllers.DeleteTokenVerification(w, r)
}
func Session_logout(w http.ResponseWriter, r *http.Request) {
	middleware.LogoutCurrentUser(w, r)

}
func LogoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.LogoutSession(w, r)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("GoogleCallback hit")
	log.Println("Cookies in request:")
	for _, c := range r.Cookies() {
		log.Printf(" - %s=%s\n", c.Name, c.Value)
	}

	r = r.WithContext(context.WithValue(r.Context(), "provider", "google"))
	controllers.GoogleCallback(w, r)
}
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "provider", "google"))

	// Print cookies sent in the request
	fmt.Println("Incoming cookies:", r.Cookies())

	// Capture the response writer to inspect headers set by Gothic
	rec := &responseRecorder{ResponseWriter: w, headers: http.Header{}}
	gothic.BeginAuthHandler(rec, r)

	// Print headers set by Gothic (including Set-Cookie)
	fmt.Println("Headers set by Gothic:", rec.Header())
}

// Struct to maintain the audit on each request, listening basically
type responseRecorder struct {
	http.ResponseWriter
	headers http.Header
}

func (r *responseRecorder) Header() http.Header {
	return r.headers
}
//Write header for the required function, this case the audit that needs the headers reconstructed
func (r *responseRecorder) WriteHeader(statusCode int) {
	for k, v := range r.headers {
		r.ResponseWriter.Header()[k] = v
	}
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}
