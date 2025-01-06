package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/AgusMolinaCode/restApi-Go.git/pkg/database"
	_ "github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

type Location struct {
	Address string  `json:"address" validate:"required"`
	Lng     float64 `json:"lng" validate:"required"`
	Lat     float64 `json:"lat" validate:"required"`
}

type DateTime struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}

type Event struct {
	ID          string              `json:"id" validate:"required,uuid4"`
	Name        string              `json:"name" validate:"required"`
	Description string              `json:"description" validate:"required"`
	Location    Location            `json:"location" validate:"required"`
	DateTimes   map[string]DateTime `json:"date_times" validate:"required,min=1"`
	UserID      string              `json:"user_id" validate:"required,uuid4"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	PaymentLink map[string]struct {
		Link  string  `json:"link"`
		Price float64 `json:"price"`
	} `json:"payment_link"`
	MinPrice         float64           `json:"min_price"`
	Tags             []string          `json:"tags"`
	TransportGuide   string            `json:"transport_guide"`
	Schedule         map[string]string `json:"schedule"`
	ExclusiveParking bool              `json:"exclusive_parking"`
	Rules            []string          `json:"rules"`
	SocialLinks      map[string]string `json:"social_links"`
	Accessibility    []string          `json:"accessibility"`
	DeliveryMethod   string            `json:"delivery_method"`
	MainImageURL     string            `json:"main_image_url"`
	AdditionalImages []string          `json:"additional_images"`
	Category         string            `json:"category"`
}

func (e Event) Save() error {
	query := `
		INSERT INTO events (id, name, description, location_address, location_lng, location_lat, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide, schedule, exclusive_parking, min_price, rules, social_links, accessibility, delivery_method, main_image_url, additional_images, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
	`

	// Convertir el mapa de DateTimes a JSON
	dateTimesJSON, err := json.Marshal(e.DateTimes)
	if err != nil {
		return err
	}

	// Convertir el mapa de Schedule a JSON
	scheduleJSON, err := json.Marshal(e.Schedule)
	if err != nil {
		return err
	}

	// Convertir el mapa de PaymentLink a JSON
	paymentLinkJSON, err := json.Marshal(e.PaymentLink)
	if err != nil {
		return err
	}

	// Convertir otros mapas a JSON
	rulesJSON, err := json.Marshal(e.Rules)
	if err != nil {
		return err
	}

	socialLinksJSON, err := json.Marshal(e.SocialLinks)
	if err != nil {
		return err
	}

	accessibilityJSON, err := json.Marshal(e.Accessibility)
	if err != nil {
		return err
	}

	additionalImagesJSON, err := json.Marshal(e.AdditionalImages)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(query, e.ID, e.Name, e.Description, e.Location.Address, e.Location.Lng, e.Location.Lat, dateTimesJSON, e.UserID, e.CreatedAt, e.UpdatedAt, paymentLinkJSON, pq.Array(e.Tags), e.TransportGuide, scheduleJSON, e.ExclusiveParking, e.MinPrice, rulesJSON, socialLinksJSON, accessibilityJSON, e.DeliveryMethod, e.MainImageURL, additionalImagesJSON, e.Category)
	return err
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT id, name, description, location_address, location_lng, location_lat, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide, schedule, exclusive_parking, min_price, rules, social_links, accessibility, delivery_method, main_image_url, additional_images, category FROM events`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var paymentLinkJSON, scheduleJSON, dateTimesJSON, rulesJSON, socialLinksJSON, accessibilityJSON, additionalImagesJSON []byte
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location.Address, &event.Location.Lng, &event.Location.Lat, &dateTimesJSON, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide, &scheduleJSON, &event.ExclusiveParking, &event.MinPrice, &rulesJSON, &socialLinksJSON, &accessibilityJSON, &event.DeliveryMethod, &event.MainImageURL, &additionalImagesJSON, &event.Category)
		if err != nil {
			return nil, err
		}

		// Convertir JSON a mapas
		err = json.Unmarshal(paymentLinkJSON, &event.PaymentLink)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(scheduleJSON, &event.Schedule)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(dateTimesJSON, &event.DateTimes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(rulesJSON, &event.Rules)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(socialLinksJSON, &event.SocialLinks)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(accessibilityJSON, &event.Accessibility)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(additionalImagesJSON, &event.AdditionalImages)
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
	query := `SELECT id, name, description, location_address, location_lng, location_lat, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide, schedule, exclusive_parking, min_price, rules, social_links, accessibility, delivery_method, main_image_url, additional_images, category FROM events WHERE id = $1`
	row := database.DB.QueryRow(query, id)

	var event Event
	var paymentLinkJSON, scheduleJSON, dateTimesJSON, rulesJSON, socialLinksJSON, accessibilityJSON, additionalImagesJSON []byte
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location.Address, &event.Location.Lng, &event.Location.Lat, &dateTimesJSON, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide, &scheduleJSON, &event.ExclusiveParking, &event.MinPrice, &rulesJSON, &socialLinksJSON, &accessibilityJSON, &event.DeliveryMethod, &event.MainImageURL, &additionalImagesJSON, &event.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Convertir JSON a mapas
	err = json.Unmarshal(paymentLinkJSON, &event.PaymentLink)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scheduleJSON, &event.Schedule)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dateTimesJSON, &event.DateTimes)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rulesJSON, &event.Rules)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(socialLinksJSON, &event.SocialLinks)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(accessibilityJSON, &event.Accessibility)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(additionalImagesJSON, &event.AdditionalImages)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func UpdateEventByID(id string, updatedEvent Event) error {
	query := `
		UPDATE events
		SET name = $1, description = $2, location_address = $3, location_lng = $4, location_lat = $5, date_times = $6, user_id = $7, updated_at = $8, payment_link = $9, tags = $10, transport_guide = $11, schedule = $12, exclusive_parking = $13, min_price = $14, rules = $15, social_links = $16, accessibility = $17, delivery_method = $18, main_image_url = $19, additional_images = $20, category = $21
		WHERE id = $22
	`
	// Convertir el mapa de PaymentLink, Schedule, DateTimes, Rules, SocialLinks, Accessibility y AdditionalImages a JSON para almacenarlos en la base de datos
	paymentLinkJSON, err := json.Marshal(updatedEvent.PaymentLink)
	if err != nil {
		return err
	}

	scheduleJSON, err := json.Marshal(updatedEvent.Schedule)
	if err != nil {
		return err
	}

	dateTimesJSON, err := json.Marshal(updatedEvent.DateTimes)
	if err != nil {
		return err
	}

	rulesJSON, err := json.Marshal(updatedEvent.Rules)
	if err != nil {
		return err
	}

	socialLinksJSON, err := json.Marshal(updatedEvent.SocialLinks)
	if err != nil {
		return err
	}

	accessibilityJSON, err := json.Marshal(updatedEvent.Accessibility)
	if err != nil {
		return err
	}

	additionalImagesJSON, err := json.Marshal(updatedEvent.AdditionalImages)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(query, updatedEvent.Name, updatedEvent.Description, updatedEvent.Location.Address, updatedEvent.Location.Lng, updatedEvent.Location.Lat, dateTimesJSON, updatedEvent.UserID, updatedEvent.UpdatedAt, paymentLinkJSON, pq.Array(updatedEvent.Tags), updatedEvent.TransportGuide, scheduleJSON, updatedEvent.ExclusiveParking, updatedEvent.MinPrice, rulesJSON, socialLinksJSON, accessibilityJSON, updatedEvent.DeliveryMethod, updatedEvent.MainImageURL, additionalImagesJSON, updatedEvent.Category, id)
	return err
}

func DeleteEventByID(id string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	return err
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

func GetAllCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM events`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetEventsByTags(tags []string) ([]Event, error) {
	query := `
		SELECT id, name, description, location_address, location_lng, location_lat, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide, schedule, exclusive_parking, min_price, rules, social_links, accessibility, delivery_method, main_image_url, additional_images, category
		FROM events
		WHERE tags && $1
	`
	rows, err := database.DB.Query(query, pq.Array(tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var paymentLinkJSON, scheduleJSON, dateTimesJSON, rulesJSON, socialLinksJSON, accessibilityJSON, additionalImagesJSON []byte
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location.Address, &event.Location.Lng, &event.Location.Lat, &dateTimesJSON, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide, &scheduleJSON, &event.ExclusiveParking, &event.MinPrice, &rulesJSON, &socialLinksJSON, &accessibilityJSON, &event.DeliveryMethod, &event.MainImageURL, &additionalImagesJSON, &event.Category)
		if err != nil {
			return nil, err
		}

		// Convertir JSON a mapas, manejando nulos
		if dateTimesJSON != nil {
			err = json.Unmarshal(dateTimesJSON, &event.DateTimes)
			if err != nil {
				return nil, err
			}
		}

		if paymentLinkJSON != nil {
			err = json.Unmarshal(paymentLinkJSON, &event.PaymentLink)
			if err != nil {
				return nil, err
			}
		}

		if scheduleJSON != nil {
			err = json.Unmarshal(scheduleJSON, &event.Schedule)
			if err != nil {
				return nil, err
			}
		}

		if rulesJSON != nil {
			err = json.Unmarshal(rulesJSON, &event.Rules)
			if err != nil {
				return nil, err
			}
		}

		if socialLinksJSON != nil {
			err = json.Unmarshal(socialLinksJSON, &event.SocialLinks)
			if err != nil {
				return nil, err
			}
		}

		if accessibilityJSON != nil {
			err = json.Unmarshal(accessibilityJSON, &event.Accessibility)
			if err != nil {
				return nil, err
			}
		}

		if additionalImagesJSON != nil {
			err = json.Unmarshal(additionalImagesJSON, &event.AdditionalImages)
			if err != nil {
				return nil, err
			}
		}

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventsByCategory(category string) ([]Event, error) {
	query := `
		SELECT id, name, description, location_address, location_lng, location_lat, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide, schedule, exclusive_parking, min_price, rules, social_links, accessibility, delivery_method, main_image_url, additional_images, category
		FROM events
		WHERE category = $1
	`
	rows, err := database.DB.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var paymentLinkJSON, scheduleJSON, dateTimesJSON, rulesJSON, socialLinksJSON, accessibilityJSON, additionalImagesJSON []byte
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location.Address, &event.Location.Lng, &event.Location.Lat, &dateTimesJSON, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide, &scheduleJSON, &event.ExclusiveParking, &event.MinPrice, &rulesJSON, &socialLinksJSON, &accessibilityJSON, &event.DeliveryMethod, &event.MainImageURL, &additionalImagesJSON, &event.Category)
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

func GetEventsByDate(date string) ([]Event, error) {
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return nil, err
	}
	formattedDate := parsedDate.Format("02/01/2006")

	query := `
		SELECT id, name, description, location_address, location_lng, location_lat, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide, schedule, exclusive_parking, min_price, rules, social_links, accessibility, delivery_method, main_image_url, additional_images, category
		FROM events
		WHERE date_times ? $1
	`
	rows, err := database.DB.Query(query, formattedDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var paymentLinkJSON, scheduleJSON, dateTimesJSON, rulesJSON, socialLinksJSON, accessibilityJSON, additionalImagesJSON []byte
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location.Address, &event.Location.Lng, &event.Location.Lat, &dateTimesJSON, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide, &scheduleJSON, &event.ExclusiveParking, &event.MinPrice, &rulesJSON, &socialLinksJSON, &accessibilityJSON, &event.DeliveryMethod, &event.MainImageURL, &additionalImagesJSON, &event.Category)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(paymentLinkJSON, &event.PaymentLink)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(scheduleJSON, &event.Schedule)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(dateTimesJSON, &event.DateTimes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(rulesJSON, &event.Rules)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(socialLinksJSON, &event.SocialLinks)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(accessibilityJSON, &event.Accessibility)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(additionalImagesJSON, &event.AdditionalImages)
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

func GetEventsByName(name string) ([]Event, error) {
	query := `
		SELECT id, name, description, location_address, location_lng, location_lat, date_times, user_id, created_at, updated_at, payment_link, tags, transport_guide, schedule, exclusive_parking, min_price, rules, social_links, accessibility, delivery_method, main_image_url, additional_images, category
		FROM events
		WHERE LOWER(name) LIKE LOWER($1 || '%')
	`
	rows, err := database.DB.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var paymentLinkJSON, scheduleJSON, dateTimesJSON, rulesJSON, socialLinksJSON, accessibilityJSON, additionalImagesJSON []byte
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location.Address, &event.Location.Lng, &event.Location.Lat, &dateTimesJSON, &event.UserID, &event.CreatedAt, &event.UpdatedAt, &paymentLinkJSON, pq.Array(&event.Tags), &event.TransportGuide, &scheduleJSON, &event.ExclusiveParking, &event.MinPrice, &rulesJSON, &socialLinksJSON, &accessibilityJSON, &event.DeliveryMethod, &event.MainImageURL, &additionalImagesJSON, &event.Category)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(paymentLinkJSON, &event.PaymentLink)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(scheduleJSON, &event.Schedule)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(dateTimesJSON, &event.DateTimes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(rulesJSON, &event.Rules)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(socialLinksJSON, &event.SocialLinks)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(accessibilityJSON, &event.Accessibility)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(additionalImagesJSON, &event.AdditionalImages)
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

func GetEventSummaries(page, limit int) ([]EventSummary, error) {
	offset := (page - 1) * limit // Calcular el desplazamiento

	query := `
		SELECT name, main_image_url, date_times, min_price
		FROM events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := database.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []EventSummary
	for rows.Next() {
		var summary EventSummary
		var dateTimesJSON []byte
		err := rows.Scan(&summary.Name, &summary.MainImageURL, &dateTimesJSON, &summary.MinPrice)
		if err != nil {
			return nil, err
		}

		var dateTimes map[string]DateTime
		err = json.Unmarshal(dateTimesJSON, &dateTimes)
		if err != nil {
			return nil, err
		}

		// Obtener la primera fecha disponible
		for date, dateTime := range dateTimes {
			if dateTime.Status == "disponibles" || dateTime.Status == "pocas unidades" {
				summary.FirstAvailableDate = date
				break
			}
		}

		summaries = append(summaries, summary)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}

type EventSummary struct {
	Name               string  `json:"name"`
	FirstAvailableDate string  `json:"first_available_date"`
	MainImageURL       string  `json:"main_image_url"`
	MinPrice           float64 `json:"min_price"`
}
