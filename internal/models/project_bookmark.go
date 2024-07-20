package models

import (
	"time"

	"github.com/priyam-trambadia/project-wall/internal/logger"
)

type ProjectBookmark struct {
	ProjectID int64     `json:"project_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (projectBookmark *ProjectBookmark) Insert() error {
	logger := logger.Logger{Caller: "ProjectBookmark::Insert model"}

	query := ` 
		INSERT INTO project_bookmarks (project_id, user_id)
		VALUES ($1, $2);
	`
	_, err := database.Exec(query, projectBookmark.ProjectID, projectBookmark.UserID)
	err = logger.AppendError(err)

	return err
}

func (projectBookmark *ProjectBookmark) Get() error {
	logger := logger.Logger{Caller: "ProjectBookmark::Get model"}

	query := `
		SELECT created_at 
		FROM project_bookmarks
		WHERE project_id = $1 AND user_id = $2;
	`
	args := []interface{}{
		projectBookmark.ProjectID,
		projectBookmark.UserID,
	}

	err := database.QueryRow(query, args...).Scan(&projectBookmark.CreatedAt)
	err = logger.AppendError(err)

	return err
}

func (projectBookmark *ProjectBookmark) Delete() error {
	logger := logger.Logger{Caller: "ProjectBookmark::Delete model"}

	query := ` 
		DELETE FROM project_bookmarks
		WHERE project_id = $1 AND user_id = $2;
	`
	_, err := database.Exec(query, projectBookmark.ProjectID, projectBookmark.UserID)
	err = logger.AppendError(err)
	return err
}

func (projectBookmark *ProjectBookmark) GetUserBookmarkStatus() (bool, error) {
	logger := logger.Logger{Caller: "GetUserBookmarkStatus model"}

	query := `
		SELECT EXISTS (
			SELECT *
			FROM project_bookmarks
			WHERE user_id = $1 AND project_id = $2
		);
	`
	args := []interface{}{
		projectBookmark.UserID,
		projectBookmark.ProjectID,
	}

	var status bool
	if err := database.QueryRow(query, args...).Scan(&status); err != nil {
		err = logger.AppendError(err)
		return false, err
	}

	return status, nil
}

func GetUserBookmarkedProjectIDs(userID int64) ([]int64, error) {
	logger := logger.Logger{Caller: "GetUserBookmarkedProjectIDs model"}

	query := `
		SELECT project_id
		FROM project_bookmarks
		WHERE user_id = $1;
	`
	rows, err := database.Query(query, userID)

	if err != nil {
		err = logger.AppendError(err)
		return nil, err
	}

	defer rows.Close()
	projects := make([]int64, 0)

	for rows.Next() {
		var projectID int64

		if err := rows.Scan(&projectID); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}
		projects = append(projects, projectID)
	}

	return projects, nil
}

func GetProjectBookmarkCount(projectID int64) (int64, error) {
	logger := logger.Logger{Caller: "GetProjectBookmarkCount model"}

	query := `
		SELECT COUNT(user_id)
		FROM project_bookmarks
		WHERE project_id = $1;
	`
	var count int64
	if err := database.QueryRow(query, projectID).Scan(&count); err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return count, nil
}
