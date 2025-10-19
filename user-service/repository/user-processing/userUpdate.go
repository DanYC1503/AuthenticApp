package repository

import (
	"database/sql"
	"fmt"
	"main/models"
	"net/http"
	"strings"
)

func UpdateUser(w http.ResponseWriter, user models.UserUpdate, tx *sql.Tx) error {

	// Build dynamic UPDATE query
	var setParts []string
	var args []interface{}
	argIndex := 1

	if user.IDNumberEncrypted != "" {
		setParts = append(setParts, fmt.Sprintf("id_number_encrypted = $%d", argIndex))
		args = append(args, user.IDNumberEncrypted)
		argIndex++
	}
	if user.FullName != "" {
		setParts = append(setParts, fmt.Sprintf("full_name = $%d", argIndex))
		args = append(args, user.FullName)
		argIndex++
	}
	if user.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, user.Email)
		argIndex++
	}
	if user.PhoneNumber != "" {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", argIndex))
		args = append(args, user.PhoneNumber)
		argIndex++
	}
	if !user.DateOfBirth.ToTime().IsZero() {
		setParts = append(setParts, fmt.Sprintf("date_of_birth = $%d", argIndex))
		args = append(args, user.DateOfBirth)
		argIndex++
	}
	if user.Address != "" {
		setParts = append(setParts, fmt.Sprintf("address = $%d", argIndex))
		args = append(args, user.Address)
		argIndex++
	}
	if user.UsernameNew != "" && user.UsernameNew != user.Username {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, user.UsernameNew)
		argIndex++
	}

	// Validate
	if len(setParts) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return fmt.Errorf("no fields to update")
	}

	// Build final query
	query := fmt.Sprintf(
		`UPDATE users SET %s WHERE username = $%d`,
		strings.Join(setParts, ", "),
		argIndex,
	)
	args = append(args, user.Username)

	// Execute the query
	result, err := tx.Exec(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found or no changes made", http.StatusNotFound)
		return fmt.Errorf("no rows updated")
	}

	return nil
}
