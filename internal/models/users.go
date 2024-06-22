package models

import (
	"time"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	Password     string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) Insert() {
	query := `
						INSERT INTO users (name, email, password)
						VALUES ($1, $2, $3);
					 `

	args := []interface{}{
		u.Name,
		u.Email,
		u.Password,
	}

	database.QueryRow(query, args...)
}

func (u *User) ValidateUser() bool {
	query := `
						SELECT _id
						FROM users
						WHERE email = $1 AND password = $2;
					 `

	args := []interface{}{
		u.Email,
		u.Password,
	}

	err := database.QueryRow(query, args...).Scan(&u.ID)

	if err != nil {
		return false
	} else {
		return true
	}
}

func (u *User) UpdateRefreshToken() {

	query := `
						UPDATE users
						SET refresh_token = $2
						WHERE _id = $1;
					 `

	args := []interface{}{
		u.ID,
		u.RefreshToken,
	}

	database.QueryRow(query, args...)
}

func (u *User) Delete() {

}
