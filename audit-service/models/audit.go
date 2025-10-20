package models

import "time"

type AuditLog struct {
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Metadata  string    `json:"metadata,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
