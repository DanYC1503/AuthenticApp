package models

type UserAuditLogs struct {
	Email     string `json:"email"`
	Action    string `json:"action"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent"`
	Metadata  string `json:"metadata,omitempty"`
	Timestamp string `json:"audit_time,omitempty"` // match the SQL column
}
type UserAdminRetrieval struct {
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	DateOfBirth   string `json:"date_of_birth"` // Change to string
	Address       string `json:"address"`
	CreateDate    string `json:"create_date"` // Change to string
	AccountStatus string `json:"account_status"`
	OAuthProvider string `json:"oauth_provider,omitempty"`
	OAuthID       string `json:"oauth_id,omitempty"`
	UserType      string `json:"user_type"`
}
