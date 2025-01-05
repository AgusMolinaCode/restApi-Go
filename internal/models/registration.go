package models

import (
	"database/sql"

	"github.com/AgusMolinaCode/restApi-Go.git/pkg/database"
)

type Registration struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	UserID      string `json:"user_id"`
	Whatsapp    string `json:"whatsapp"`
	CreatedAt   string `json:"created_at"`
	EventDate   string `json:"event_date"`
	PaymentLink string `json:"payment_link"`
}

func (r *Registration) Save() error {
	query := `INSERT INTO registrations (id, event_id, user_id, whatsapp, created_at, event_date, payment_link) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.ID, r.EventID, r.UserID, r.Whatsapp, r.CreatedAt, r.EventDate, r.PaymentLink)
	return err
}

func IsUserRegisteredForEvent(eventID, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM registrations WHERE event_id = $1 AND user_id = $2`
	row := database.DB.QueryRow(query, eventID, userID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteRegistration(eventID, userID string) error {
	query := `DELETE FROM registrations WHERE event_id = $1 AND user_id = $2`
	_, err := database.DB.Exec(query, eventID, userID)
	return err
}

type RegistrationDetail struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Whatsapp  string `json:"whatsapp"`
	CreatedAt string `json:"created_at"`
}

func GetRegistrationsByEventID(eventID string) ([]RegistrationDetail, error) {
	query := `
		SELECT users.id, users.username, users.email, users.whatsapp, registrations.created_at
		FROM registrations
		JOIN users ON registrations.user_id = users.id
		WHERE registrations.event_id = $1
	`
	rows, err := database.DB.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []RegistrationDetail
	for rows.Next() {
		var reg RegistrationDetail
		err := rows.Scan(&reg.UserID, &reg.Username, &reg.Email, &reg.Whatsapp, &reg.CreatedAt)
		if err != nil {
			return nil, err
		}
		registrations = append(registrations, reg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

func GetRegistrationByUserID(eventID, userID string) (*RegistrationDetail, error) {
	query := `
		SELECT users.id, users.username, users.email, users.whatsapp, registrations.created_at
		FROM registrations
		JOIN users ON registrations.user_id = users.id
		WHERE registrations.event_id = $1 AND registrations.user_id = $2
	`
	row := database.DB.QueryRow(query, eventID, userID)

	var reg RegistrationDetail
	err := row.Scan(&reg.UserID, &reg.Username, &reg.Email, &reg.Whatsapp, &reg.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reg, nil
}
