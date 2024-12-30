package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/AgusMolinaCode/restApi-Go.git/pkg/database"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id" validate:"required,uuid4"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *User) Save() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	fmt.Println("Hashed Password:", u.Password)

	query := `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
	`
	_, err = database.DB.Exec(query, u.ID, u.Username, u.Email, u.Password)
	return err
}

func GetUserByID(id string) (*UserResponse, error) {
	query := `SELECT id, username, email FROM users WHERE id = $1`
	row := database.DB.QueryRow(query, id)

	var user UserResponse
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]UserResponse, error) {
	query := `SELECT id, username, email FROM users`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var user UserResponse
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
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

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	row := database.DB.QueryRow(query, email)

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

func UpdateUserByID(id string, updatedUser User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	updatedUser.Password = string(hashedPassword)

	query := `UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4`
	_, err = database.DB.Exec(query, updatedUser.Username, updatedUser.Email, updatedUser.Password, id)
	return err
}

func DeleteUserByID(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	return err
}

func SetResetToken(email string) (string, error) {
	token := uuid.New().String()
	expiry := time.Now().Add(1 * time.Hour)

	query := `UPDATE users SET reset_token = $1, reset_token_expiry = $2 WHERE email = $3`
	_, err := database.DB.Exec(query, token, expiry, email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyResetToken(token string) (string, error) {
	query := `SELECT id FROM users WHERE reset_token = $1 AND reset_token_expiry > $2`
	row := database.DB.QueryRow(query, token, time.Now())

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return userID, nil
}

func UpdatePassword(userID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = $1, reset_token = NULL, reset_token_expiry = NULL WHERE id = $2`
	_, err = database.DB.Exec(query, string(hashedPassword), userID)
	return err
}
