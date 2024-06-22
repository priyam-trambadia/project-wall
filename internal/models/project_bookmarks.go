package models

import "time"

type ProjectBookmark struct {
	ProjectID int64
	UserID    int64
	CreatedAt time.Time
}

func (projectBookmark *ProjectBookmark) Insert() {
	query := ` 
		INSERT INTO project_bookmarks (project_id, user_id)
		VALUES ($1, $2);
	`
	database.QueryRow(query, projectBookmark.ProjectID, projectBookmark.UserID)
}

func (projectBookmark *ProjectBookmark) Delete() {
	query := ` 
		DELETE FROM project_bookmarks
		WHERE project_id = $1 AND user_id = $2;
	`

	database.QueryRow(query, projectBookmark.ProjectID, projectBookmark.UserID)
}

func (projectBookmark *ProjectBookmark) GetBookmarks() []int64 {
	bookmarks := make([]int64, 0)

	query := `
		SELECT user_id
		FROM project_bookmarks
		WHERE project_id = $1;
	`

	rows, _ := database.Query(query, projectBookmark.ProjectID)
	defer rows.Close()

	for rows.Next() {
		var UserID int64

		err := rows.Scan(&UserID)
		if err != nil {
			break
		}

		bookmarks = append(bookmarks, UserID)
	}

	return bookmarks
}

func (projectBookmark *ProjectBookmark) GetProjects() []int64 {
	projects := make([]int64, 0)

	query := `
		SELECT project_id
		FROM project_bookmarks
		WHERE user_id = $1;
	`
	rows, _ := database.Query(query, projectBookmark.UserID)
	defer rows.Close()

	for rows.Next() {
		var projectID int64

		err := rows.Scan(&projectID)
		if err != nil {
			break
		}

		projects = append(projects, projectID)
	}

	return projects
}
