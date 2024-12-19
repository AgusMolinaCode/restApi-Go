package models

import (
	"time"
	"github.com/go-playground/validator/v10"
)

type Event struct {
	ID          string    `json:"id" validate:"required,uuid4"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	DateTime    time.Time `json:"date_time" validate:"required"`
	UserID      string    `json:"user_id" validate:"required,uuid4"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var events = []Event{}

func (e Event) Save() error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	// TODO: Implementar la lógica de guardado del evento y añadir la base de datos
	events = append(events, e)
	return nil
}

func GetAllEvents() []Event {
	return events
}


