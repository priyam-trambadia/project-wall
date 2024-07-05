package models

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID             int64
	Name           string
	Email          string
	Password       string
	Bio            string
	SocialLinks    []string
	LocationID     int64
	OrganizationID int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RefreshToken   string
}

func (user *User) Insert() error {
	query := `
		INSERT INTO users (
			name, 
			email, 
			password, 
			organization_id
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	parts := strings.Split(user.Email, "@")
	if len := len(parts); len < 2 {
		return errors.New("invalid email format")
	} else {
		hostname := parts[len-1]

		var err error
		if user.OrganizationID, err = GetOrCreateOrganizationID(hostname); err != nil {
			return err
		}
	}

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password,
		user.OrganizationID,
	}

	return database.QueryRow(query, args...).Scan(&user.ID)
}

func (user *User) Update() error {
	query := `
		UPDATE users 
		SET 
			name = $1,
			bio = $2,
			social_links = $3,
			location_id = $4,
			updated_at = $5
		WHERE id = $6;
	`
	args := []interface{}{
		user.Name,
		user.Bio,
		user.SocialLinks,
		user.LocationID,
		time.Now(),
		user.ID,
	}

	_, err := database.Exec(query, args...)
	return err
}

func (user *User) Get() error {
	query := `
		SELECT 
			name,
			email,
			password,
			bio,
			social_links,
			location_id,
			organization_id,
			created_at,
			updated_at,
			refresh_token
		FROM users
		WHERE id = $1;
	`

	err := database.QueryRow(query, user.ID).Scan(
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.SocialLinks,
		&user.LocationID,
		&user.OrganizationID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.RefreshToken,
	)
	return err
}

func (user *User) Delete() error {
	query := `
		DELETE FROM users
		WHERE id = $1;
	`
	_, err := database.Exec(query, user.ID)
	return err
}

// utilities

func GetUserID(email, password string) (int64, error) {
	query := `
		SELECT id
		FROM users
		WHERE email = $1 AND password = $2;
	`
	args := []interface{}{
		email,
		password,
	}

	var id int64
	err := database.QueryRow(query, args...).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateUserRefreshToken(id int64, refresh_token string) error {

	query := `
		UPDATE users
		SET refresh_token = $1
		WHERE id = $2;
	`

	args := []interface{}{
		id,
		refresh_token,
	}

	_, err := database.Exec(query, args...)
	return err
}

func IsUserExists(userID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		);
	`
	var exists bool
	if err := database.QueryRow(query, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
