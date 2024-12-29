package models

import (
	"database/sql"

	"github.com/AgusMolinaCode/restApi-Go.git/pkg/database"
	_ "github.com/go-playground/validator/v10"
)

type Event struct {
	ID          string `json:"id" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Location    string `json:"location" validate:"required"`
	DateTime    string `json:"date_time" validate:"required"`
	UserID      string `json:"user_id" validate:"required,uuid4"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	PaymentLink string `json:"payment_link"`
}

func (e Event) Save() error {
	query := `INSERT INTO events (id, name, description, location, date_time, user_id, created_at, updated_at, payment_link) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, e.Name, e.Description, e.Location, e.DateTime, e.UserID, e.CreatedAt, e.UpdatedAt, e.PaymentLink)
	if err != nil {
		return err
	}
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &event.PaymentLink)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventByID(id string) (*Event, error) {
	query := `SELECT id, name, description, location, date_time, user_id, created_at, updated_at, payment_link FROM events WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &event.PaymentLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func UpdateEventByID(id string, updatedEvent Event) error {
	query := `UPDATE events SET name = ?, description = ?, location = ?, date_time = ?, user_id = ?, updated_at = ?, payment_link = ? WHERE id = ?`
	_, err := database.DB.Exec(query, updatedEvent.Name, updatedEvent.Description, updatedEvent.Location, updatedEvent.DateTime, updatedEvent.UserID, updatedEvent.UpdatedAt, updatedEvent.PaymentLink, id)
	return err
}

func DeleteEventByID(id string) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}

type Registration struct {
	ID        string `json:"id"`
	EventID   string `json:"event_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (r *Registration) Save() error {
	query := `INSERT INTO registrations (id, event_id, user_id, created_at) VALUES (?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.ID, r.EventID, r.UserID, r.CreatedAt)
	return err
}

func IsUserRegisteredForEvent(eventID, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM registrations WHERE event_id = ? AND user_id = ?`
	row := database.DB.QueryRow(query, eventID, userID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteRegistration(eventID, userID string) error {
	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, eventID, userID)
	return err
}

type RegistrationDetail struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func GetRegistrationsByEventID(eventID string) ([]RegistrationDetail, error) {
	query := `
		SELECT users.id, users.username, users.email, registrations.created_at
		FROM registrations
		JOIN users ON registrations.user_id = users.id
		WHERE registrations.event_id = ?
	`
	rows, err := database.DB.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []RegistrationDetail
	for rows.Next() {
		var reg RegistrationDetail
		err := rows.Scan(&reg.UserID, &reg.Username, &reg.Email, &reg.CreatedAt)
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
