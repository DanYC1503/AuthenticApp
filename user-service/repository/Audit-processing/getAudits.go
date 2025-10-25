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
		var phoneNumber, dateOfBirth, address sql.NullString
		var createDate sql.NullString
		var accountStatus sql.NullString

		if err := rows.Scan(
			&u.Username,
			&u.FullName,
			&u.Email,
			&phoneNumber,
			&dateOfBirth,
			&address,
			&createDate,
			&accountStatus,
			&u.OAuthProvider,
			&u.OAuthID,
			&u.UserType,
		); err != nil {
			return nil, err
		}

		// Convert nullable fields to strings safely
		u.PhoneNumber = phoneNumber.String
		u.DateOfBirth = dateOfBirth.String
		u.Address = address.String
		u.CreateDate = createDate.String
		u.AccountStatus = accountStatus.String

		users = append(users, u)
	}

	return users, nil
}
