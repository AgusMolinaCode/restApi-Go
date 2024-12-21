package models

import (
	"database/sql"

	"github.com/AgusMolinaCode/restApi-Go.git/db"
	_ "github.com/go-playground/validator/v10"
)

type User struct {
	ID       string `json:"id" validate:"required,uuid4"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u User) Save() error {
	query := `INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.ID, u.Username, u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(id string) (*User, error) {
	query := `SELECT id, username, email, password FROM users WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id, username, email, password FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
