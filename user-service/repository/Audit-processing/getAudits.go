package repository

import (
	"database/sql"
	"main/models"
)

func UserAuditActions(tx *sql.Tx, email string) ([]models.UserAuditLogs, error) {
	rows, err := tx.Query(`SELECT * FROM retrieve_user_audits($1)`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []models.UserAuditLogs
	for rows.Next() {
		var a models.UserAuditLogs
		if err := rows.Scan(&a.Email, &a.Action, &a.IPAddress, &a.UserAgent, &a.Metadata, &a.Timestamp); err != nil {
			return nil, err
		}
		audits = append(audits, a)
	}

	return audits, nil
}
func GetUsers(db *sql.DB, email string) ([]models.UserAdminRetrieval, error) {
	rows, err := db.Query(`SELECT * FROM retrieve_users($1)`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserAdminRetrieval
	for rows.Next() {
		var u models.UserAdminRetrieval
		if err := rows.Scan(
			&u.Username,
			&u.FullName,
			&u.Email,
			&u.PhoneNumber,
			&u.DateOfBirth,
			&u.Address,
			&u.CreateDate,
			&u.AccountStatus,
			&u.OAuthProvider,
			&u.OAuthID,
			&u.UserType,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
