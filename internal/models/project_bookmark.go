package models

import "time"

type ProjectBookmark struct {
	ProjectID int64
	UserID    int64
	CreatedAt time.Time
}

func (projectBookmark *ProjectBookmark) Insert() error {
	query := ` 
		INSERT INTO project_bookmarks (project_id, user_id)
		VALUES ($1, $2);
	`
	_, err := database.Exec(query, projectBookmark.ProjectID, projectBookmark.UserID)
	return err
}

func (projectBookmark *ProjectBookmark) Get() error {
	query := `
		SELECT created_at 
		FROM project_bookmarks
		WHERE project_id = $1 AND user_id = $2;
	`
	args := []interface{}{
		projectBookmark.ProjectID,
		projectBookmark.UserID,
	}

	return database.QueryRow(query, args...).Scan(&projectBookmark.CreatedAt)
}

func (projectBookmark *ProjectBookmark) Delete() error {
	query := ` 
		DELETE FROM project_bookmarks
		WHERE project_id = $1 AND user_id = $2;
	`
	_, err := database.Exec(query, projectBookmark.ProjectID, projectBookmark.UserID)
	return err
}

func GetUserBookmarkedProjectIDs(userID int64) ([]int64, error) {
	query := `
		SELECT project_id
		FROM project_bookmarks
		WHERE user_id = $1;
	`
	if rows, err := database.Query(query, userID); err != nil {
		return nil, err
	} else {

		defer rows.Close()
		projects := make([]int64, 0)

		for rows.Next() {
			var projectID int64

			if err := rows.Scan(&projectID); err != nil {
				return nil, err
			} else {
				projects = append(projects, projectID)
			}
		}
		return projects, nil
	}
}
