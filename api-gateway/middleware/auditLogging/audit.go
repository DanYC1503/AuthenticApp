package auditlogging

type AuditEvent struct {
	Username   string
	IP         string
	UA         string
	Method     string
	Path       string
	StatusCode int
}
