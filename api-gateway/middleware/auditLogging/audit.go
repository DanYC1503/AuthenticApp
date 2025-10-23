package auditlogging

type auditEvent struct {
	username string
	ip       string
	ua       string
	method   string
	path     string
}
