package models

import (
	"database/sql"
	"github.com/AgusMolinaCode/restApi-Go.git/db"
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
}

func (e Event) Save() error {
	query := `INSERT INTO events (id, name, description, location, date_time, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, e.Name, e.Description, e.Location, e.DateTime, e.UserID, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID, &event.CreatedAt, &event.UpdatedAt)
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
	query := `SELECT id, name, description, location, date_time, user_id, created_at, updated_at FROM events WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func UpdateEventByID(id string, updatedEvent Event) error {
	query := `UPDATE events SET name = ?, description = ?, location = ?, date_time = ?, user_id = ?, updated_at = ? WHERE id = ?`
	_, err := db.DB.Exec(query, updatedEvent.Name, updatedEvent.Description, updatedEvent.Location, updatedEvent.DateTime, updatedEvent.UserID, updatedEvent.UpdatedAt, id)
	return err
}

func DeleteEventByID(id string) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	return err
}
