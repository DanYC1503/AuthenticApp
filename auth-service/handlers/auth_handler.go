package handlers

import (
	"context"
	"fmt"
	"log"
	service "main/controllers"
	"main/middleware"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func GetDeleteToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	service.GetDeleteToken(w, r)
}
func GetUpdateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	service.GetUpdateToken(w, r)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	service.CreateUser(w, r)
}
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	service.LoginUser(w, r)
}
func Session_logout(w http.ResponseWriter, r *http.Request) {
	middleware.LogoutCurrentUser(w, r)

}
func TokenVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	service.SessionTokenVerification(w, r)
}
func UpTokenVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	service.UpdateTokenVerification(w, r)
}
func DelTokenVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	service.DeleteTokenVerification(w, r)
}
func LogoutSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	service.LogoutSession(w, r)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("GoogleCallback hit")
	log.Println("Cookies in request:")
	for _, c := range r.Cookies() {
		log.Printf(" - %s=%s\n", c.Name, c.Value)
	}

	r = r.WithContext(context.WithValue(r.Context(), "provider", "google"))
	service.GoogleCallback(w, r)
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

type responseRecorder struct {
	http.ResponseWriter
	headers http.Header
}

func (r *responseRecorder) Header() http.Header {
	return r.headers
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	for k, v := range r.headers {
		r.ResponseWriter.Header()[k] = v
	}
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}
