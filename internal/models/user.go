package models

import (
	"errors"
	"strings"
	"time"

	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
	"github.com/priyam-trambadia/project-wall/internal/logger"
)

type User struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Avatar         string    `json:"avatar"`
	Bio            string    `json:"bio"`
	SocialLinks    string    `json:"social_links"`
	OrganizationID int64     `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	RefreshToken   string    `json:"refresh_token"`
	IsActivated    bool      `json:"is_activated"`
}

func (user *User) Insert() error {
	logger := logger.Logger{Caller: "User::Insert model"}

	query := `
		INSERT INTO users (
			name, 
			email, 
			password, 
			organization_id,
			refresh_token
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	parts := strings.Split(user.Email, "@")
	len := len(parts)
	if len < 2 {
		err := errors.New("invalid email format")
		err = logger.AppendError(err)
		return err
	}

	hostname := parts[len-1]

	organizationID, err := GetOrCreateOrganizationID(hostname)
	if err != nil {
		err = logger.AppendError(err)
		return err
	}
	user.OrganizationID = organizationID

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password,
		user.OrganizationID,
		jwt.ValEmptyToken,
	}

	err2 := database.QueryRow(query, args...).Scan(&user.ID)
	if err2 != nil {
		err2 = logger.AppendError(err2)
	}
	return err2
}

func (user *User) Update() error {
	logger := logger.Logger{Caller: "User::Update model"}

	query := `
		UPDATE users 
		SET 
			name = $1,
			bio = $2,
			social_links = $3,
			updated_at = $4
		WHERE id = $5;
	`
	args := []interface{}{
		user.Name,
		user.Bio,
		user.SocialLinks,
		time.Now(),
		user.ID,
	}

	_, err := database.Exec(query, args...)
	err = logger.AppendError(err)
	return err
}

func (user *User) Get() error {
	logger := logger.Logger{Caller: "User::Get model"}

	query := `
		SELECT 
			name,
			email,
			password,
			avatar,
			bio,
			social_links,
			organization_id,
			created_at,
			updated_at,
			refresh_token,
			is_activated 
		FROM users
		WHERE id = $1;
	`

	err := database.QueryRow(query, user.ID).Scan(
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Bio,
		&user.SocialLinks,
		&user.OrganizationID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.RefreshToken,
		&user.IsActivated,
	)

	err = logger.AppendError(err)
	return err
}

func (user *User) Delete() error {
	logger := logger.Logger{Caller: "User::Delete model"}

	query := `
		DELETE FROM users
		WHERE id = $1;
	`
	_, err := database.Exec(query, user.ID)
	err = logger.AppendError(err)
	return err
}

// utilities

func IsEmailExists(email string) (bool, error) {
	logger := logger.Logger{Caller: "IsEmailExists model"}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		);
	`

	var exists bool
	if err := database.QueryRow(query, email).Scan(&exists); err != nil {
		err = logger.AppendError(err)
		return false, err
	}

	return exists, nil
}

func GetUserID(email, password string) (int64, error) {
	logger := logger.Logger{Caller: "GetUserID model"}

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
		err = logger.AppendError(err)
		return 0, err
	}

	return id, nil
}

func GetUserIDByEmail(email string) (int64, error) {
	logger := logger.Logger{Caller: "GetUserIDByEmail model"}

	query := `
		SELECT id
		FROM users
		WHERE email = $1;
	`
	var id int64
	err := database.QueryRow(query, email).Scan(&id)
	if err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return id, nil
}

func GetUserOrganizationID(userID int64) (int64, error) {
	logger := logger.Logger{Caller: "GetUserOrganizationID model"}

	query := `
		SELECT organization_id
		FROM users
		WHERE id = $1;
	`

	var organizationID int64
	err := database.QueryRow(query, userID).Scan(&organizationID)
	if err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return organizationID, nil
}

func UpdateUserRefreshToken(id int64, refresh_token string) error {
	logger := logger.Logger{Caller: "UpdateUserRefreshToken model"}

	query := `
		UPDATE users
		SET refresh_token = $1
		WHERE id = $2;
	`

	args := []interface{}{
		refresh_token,
		id,
	}

	_, err := database.Exec(query, args...)
	err = logger.AppendError(err)
	return err
}

func IsUserExists(userID int64) (bool, error) {
	logger := logger.Logger{Caller: "IsUserExists model"}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		);
	`
	var exists bool
	if err := database.QueryRow(query, userID).Scan(&exists); err != nil {
		err = logger.AppendError(err)
		return false, err
	}

	return exists, nil
}

func ActivateUser(userID int64) error {
	logger := logger.Logger{Caller: "ActivateUser model"}

	query := `
		UPDATE users 
		SET is_activated = $1
		WHERE id = $2;
	`
	args := []interface{}{
		true,
		userID,
	}
	_, err := database.Exec(query, args...)
	err = logger.AppendError(err)
	return err
}

func UpdateUserPassword(userID int64, newPassword string) error {
	logger := logger.Logger{Caller: "UpdateUserPassword model"}

	query := `
	UPDATE users 
	SET password = $1
	WHERE id = $2;
`
	args := []interface{}{
		newPassword,
		userID,
	}
	_, err := database.Exec(query, args...)
	err = logger.AppendError(err)
	return err
}
