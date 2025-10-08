package models

import (
	"errors"

	"github.com/lilbonekit/event-management-svg/db"
	"github.com/lilbonekit/event-management-svg/utils"
)

type User struct {
	ID       int64  `json:"id" example:"1"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"mypassword123"`
}

func (u *User) Save() error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		panic("Could not prepare statement.")
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	return nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		return err
	}

	isPasswordValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !isPasswordValid {
		return errors.New("invalid credentials")
	}

	return nil
}
