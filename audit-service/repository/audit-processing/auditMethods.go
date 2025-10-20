package auditprocessing

import (
	"database/sql"
	"main/models"

	_ "github.com/lib/pq"
)

func AuditAction(tx *sql.Tx, audit models.AuditLog) error {
	query := `SELECT log_user_action($1, $2, $3, $4, $5)`

	// Just execute and ignore the result
	_, err := tx.Exec(query,
		audit.Username,
		audit.Action,
		audit.IPAddress,
		audit.UserAgent,
		audit.Metadata,
	)

	return err
}
