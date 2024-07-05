package models

import (
	"fmt"
	"time"
)

type Project struct {
	ID          int64     `json:"id"`
	GithubURL   string    `json:"github_url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerID     int64     `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags"`
	Languages   []string  `json:"languages"`
}

func (project *Project) Insert() error {
	query := `
		INSERT INTO projects (github_url, title, description, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	args := []interface{}{
		project.GithubURL,
		project.Title,
		project.Description,
		project.OwnerID,
	}

	if err := database.QueryRow(query, args...).Scan(&project.ID); err != nil {
		return err
	} else {
		// adding tags
		for _, tagName := range project.Tags {
			if tagID, err := GetOrCreateTagID(tagName); err != nil {
				return err
			} else {
				projectTag := ProjectTag{ProjectID: project.ID, TagID: tagID}
				if err := projectTag.Insert(); err != nil {
					return err
				}
			}
		}

		//adding languages
		for _, languageName := range project.Languages {
			if languageID, err := GetOrCreateLanguageID(languageName); err != nil {
				return err
			} else {
				projectLanguage := ProjectLanguage{ProjectID: project.ID, LanguageID: languageID}
				if err := projectLanguage.Insert(); err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func (project *Project) Update() error {
	query := `
		UPDATE projects
		SET description = $1, updated_at = $2,  
		WHERE id = $3;
	`
	args := []interface{}{
		project.Description,
		time.Now(),
		project.ID,
	}

	if _, err := database.Exec(query, args...); err != nil {
		return err
	} else {
		// updating tags
		if oldTagList, err := GetProjectTagIDs(project.ID); err != nil {
			return err
		} else {
			newTagList := make([]int64, 0)
			for _, tagName := range project.Tags {
				if tagID, err := GetOrCreateTagID(tagName); err != nil {
					return err
				} else {
					newTagList = append(newTagList, tagID)
				}
			}

			if err := SyncProjectTags(project.ID, oldTagList, newTagList); err != nil {
				return err
			}
		}

		// updating languages
		if oldLanguageList, err := GetProjectLanguagesIDs(project.ID); err != nil {
			return err
		} else {
			newLanguageList := make([]int64, 0)
			for _, languageName := range project.Languages {
				if languageID, err := GetOrCreateLanguageID(languageName); err != nil {
					return err
				} else {
					newLanguageList = append(newLanguageList, languageID)
				}

				if err := SyncProjectLanguages(project.ID, oldLanguageList, newLanguageList); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (project *Project) Get() error {
	query := `
		SELECT 
			github_url,
			title,
			description,
			owner_id,
			created_at,
			updated_at
		FROM projects
		WHERE id = $1;	
	`
	err := database.QueryRow(query, project.ID).Scan(
		&project.GithubURL,
		&project.Title,
		&project.OwnerID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if tagIDs, err := GetProjectTagIDs(project.ID); err != nil {
		return err
	} else {
		if project.Tags, err = GetTagNames(tagIDs); err != nil {
			return err
		}
	}

	if languageIDs, err := GetProjectLanguagesIDs(project.ID); err != nil {
		return err
	} else {
		if project.Languages, err = GetLanguageNames(languageIDs); err != nil {
			return err
		}
	}

	return nil
}

func (project *Project) Delete() error {
	query := ` 
		DELETE FROM projects
		WHERE project_id = $1;
	`

	_, err := database.Exec(query, project.ID)
	return err
}

type SortBy string

const (
	Date          SortBy = "created_at"
	BookmarkCount SortBy = "bookmark_count"
)

type SortDirection string

const (
	Ascending  SortDirection = "ASC"
	Descending SortDirection = "DESC"
)

type ProjectSearchQuery struct {
	Title          string
	TagIDs         []int64
	LanguageIDs    []int64
	OrganizationID int64
	SortBy         SortBy
	SortDirection  SortDirection
}

func (searchQuery *ProjectSearchQuery) FindProjectsWithFullTextSearch() ([]Project, error) {

	if searchQuery.SortBy == "" {
		searchQuery.SortBy = Date
	}

	if searchQuery.SortDirection == "" {
		searchQuery.SortDirection = Ascending
	}

	query := fmt.Sprintf(` 
		SELECT 
			id,
			COUNT(*) AS bookmark_count
		FROM projects AS p
		INNER JOIN project_bookmarks AS pb
		ON p.id = pb.project_id
		WHERE 
			(	$1 = '' OR
				to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) ) AND
			( $2 = 0 OR
			  id IN ( SELECT project_id
							 	FROM project_tags
					 			WHERE tag_id IN %s ) ) AND
			( $3 = 0 OR 
				id IN ( SELECT project_id
								FROM project_languages
								WHERE language_id IN %s ) ) AND
			(	$4 = 0 OR			
				owner_id IN ( SELECT id
											FROM users
											WHERE organization_id = $4 ) )
		GROUP BY p.id												
		ORDER BY %s %s;`,
		ArraytoStringRoundBrackets(searchQuery.TagIDs),
		ArraytoStringRoundBrackets(searchQuery.LanguageIDs),
		searchQuery.SortBy,
		searchQuery.SortDirection,
	)

	if rows, err := database.Query(query); err != nil {
		return nil, err
	} else {
		defer rows.Close()

		projects := make([]Project, 0)

		for rows.Next() {
			var project Project

			err2 := rows.Scan(
				&project.ID,
				&project.GithubURL,
				&project.Title,
				&project.Description,
				&project.OwnerID,
				&project.CreatedAt,
				&project.UpdatedAt,
			)

			if err2 != nil {
				return nil, err
			}

			var tagIDs []int64
			if tagIDs, err2 = GetProjectTagIDs(project.ID); err2 != nil {
				return nil, err2
			}

			if project.Tags, err2 = GetTagNames(tagIDs); err2 != nil {
				return nil, err2
			}

			var languageIDs []int64
			if languageIDs, err2 = GetProjectLanguagesIDs(project.ID); err2 != nil {
				return nil, err2
			}

			if project.Languages, err2 = GetLanguageNames(languageIDs); err2 != nil {
				return nil, err2
			}

			projects = append(projects, project)

		}
		return projects, nil
	}
}

func GetProjectOwnerID(projectID int64) (int64, error) {
	query := `
		SELECT owner_id
		FROM projects
		WHERE id = $1;
	`
	var ownerID int64
	if err := database.QueryRow(query, projectID).Scan(&ownerID); err != nil {
		return 0, err
	} else {
		return ownerID, nil
	}
}

func IsProjectExists(projectID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1	
			FROM projects
			WHERE id = $1
		);
	`
	var exists bool
	if err := database.QueryRow(query, projectID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
