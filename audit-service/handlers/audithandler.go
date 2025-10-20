package handlers

import (
	"main/controllers"
	"net/http"
)

func AuditAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	controllers.AuditActionHandler(w, r)
}
