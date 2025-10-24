package handlers

import (
	"log"
	"main/controllers"
	"net/http"
)

func AuditAction(w http.ResponseWriter, r *http.Request) {
	log.Println("[Audit] Incoming request to /audit/log")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.AuditActionHandler(w, r)
}
