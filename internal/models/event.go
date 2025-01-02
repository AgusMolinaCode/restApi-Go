package models

import (
	"database/sql"
	"encoding/json"

	"github.com/AgusMolinaCode/restApi-Go.git/pkg/database"
	_ "github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

type Event struct {
	ID             string            `json:"id" validate:"required,uuid4"`
	Name           string            `json:"name" validate:"required"`
	Description    string            `json:"description" validate:"required"`
	Location       string            `json:"location" validate:"required"`
	DateTimes      []string          `json:"date_times" validate:"required,min=1,max=2"`
	UserID         string            `json:"user_id" validate:"required,uuid4"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	PaymentLink    map[string]string `json:"payment_link"`
	Tags           []string          `json:"tags"`
	TransportGuide string            `json:"transport_guide"`
}

func (e Event) Save() error {
	query := `
		INSERT INTO events (id, name, description, location, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	// Convertir el mapa de PaymentLink a JSON para almacenarlo en la base de datos
	paymentLinkJSON, err := json.Marshal(e.PaymentLink)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(query, e.ID, e.Name, e.Description, e.Location, pq.Array(e.DateTimes), e.UserID, e.CreatedAt, e.UpdatedAt, paymentLinkJSON, pq.Array(e.Tags), e.TransportGuide)
	return err
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT id, name, description, location, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide FROM events`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var paymentLinkJSON []byte
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, pq.Array(&event.DateTimes), &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide)
		if err != nil {
			return nil, err
		}

		// Convertir JSON de PaymentLink a mapa
		err = json.Unmarshal(paymentLinkJSON, &event.PaymentLink)
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
	query := `SELECT id, name, description, location, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide FROM events WHERE id = $1`
	row := database.DB.QueryRow(query, id)

	var event Event
	var paymentLinkJSON []byte
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, pq.Array(&event.DateTimes), &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Convertir JSON de PaymentLink a mapa
	err = json.Unmarshal(paymentLinkJSON, &event.PaymentLink)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func UpdateEventByID(id string, updatedEvent Event) error {
	query := `
		UPDATE events
		SET name = $1, description = $2, location = $3, date_times = $4, user_id = $5, updated_at = $6, payment_link = $7, tags = $8, transport_guide = $9
		WHERE id = $10
	`
	// Convertir el mapa de PaymentLink a JSON para almacenarlo en la base de datos
	paymentLinkJSON, err := json.Marshal(updatedEvent.PaymentLink)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(query, updatedEvent.Name, updatedEvent.Description, updatedEvent.Location, pq.Array(updatedEvent.DateTimes), updatedEvent.UserID, updatedEvent.UpdatedAt, paymentLinkJSON, pq.Array(updatedEvent.Tags), updatedEvent.TransportGuide, id)
	return err
}

func DeleteEventByID(id string) error {
	query := `DELETE FROM events WHERE id = $1`
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
	query := `INSERT INTO registrations (id, event_id, user_id, created_at) VALUES ($1, $2, $3, $4)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.ID, r.EventID, r.UserID, r.CreatedAt)
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
	CreatedAt string `json:"created_at"`
}

func GetRegistrationsByEventID(eventID string) ([]RegistrationDetail, error) {
	query := `
		SELECT users.id, users.username, users.email, registrations.created_at
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

func GetAllTags() ([]string, error) {
	query := `SELECT DISTINCT UNNEST(tags) FROM events`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
