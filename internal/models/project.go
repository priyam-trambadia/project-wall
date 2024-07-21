package models

import (
	"fmt"
	"time"

	"github.com/priyam-trambadia/project-wall/internal/logger"
)

type Project struct {
	ID                 int64      `json:"id"`
	GithubURL          string     `json:"github_url"`
	Title              string     `json:"title"`
	Description        string     `json:"description"`
	OwnerID            int64      `json:"owner_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Tags               []Tag      `json:"tags"`
	Languages          []Language `json:"languages"`
	BookmarkCount      int64      `json:"bookmark_count"`
	UserBookmarkStatus bool       `json:"user_bookmark_status"`
}

func (project *Project) Insert() error {
	logger := logger.Logger{Caller: "Project::Insert model"}

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

	err := database.QueryRow(query, args...).Scan(&project.ID)

	if err != nil {
		err = logger.AppendError(err)
		return err
	}

	// adding tags
	for _, tag := range project.Tags {
		tagID, err := GetOrCreateTagID(tag.Name)
		if err != nil {
			err = logger.AppendError(err)
			return err
		}

		projectTag := ProjectTag{ProjectID: project.ID, TagID: tagID}
		if err := projectTag.Insert(); err != nil {
			err = logger.AppendError(err)
			return err
		}
	}

	//adding languages
	for _, language := range project.Languages {
		languageID, err := GetOrCreateLanguageID(language.Name)
		if err != nil {
			err = logger.AppendError(err)
			return err
		}

		projectLanguage := ProjectLanguage{ProjectID: project.ID, LanguageID: languageID}
		if err := projectLanguage.Insert(); err != nil {
			err = logger.AppendError(err)
			return err
		}
	}

	return nil
}

func (project *Project) Update() error {
	logger := logger.Logger{Caller: "Project::Update model"}

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

	_, err := database.Exec(query, args...)
	if err != nil {
		err = logger.AppendError(err)
		return err
	}

	// updating tags
	oldTagList, err := GetProjectTagIDs(project.ID)
	if err != nil {
		err = logger.AppendError(err)
		return err
	}

	newTagList := make([]int64, 0)
	for _, tag := range project.Tags {
		tagID, err := GetOrCreateTagID(tag.Name)
		if err != nil {
			err = logger.AppendError(err)
			return err
		}

		newTagList = append(newTagList, tagID)
	}

	if err := SyncProjectTags(project.ID, oldTagList, newTagList); err != nil {
		err = logger.AppendError(err)
		return err
	}

	// updating languages
	oldLanguageList, err := GetProjectLanguagesIDs(project.ID)
	if err != nil {
		err = logger.AppendError(err)
		return err
	}

	newLanguageList := make([]int64, 0)
	for _, language := range project.Languages {
		languageID, err := GetOrCreateLanguageID(language.Name)
		if err != nil {
			err = logger.AppendError(err)
			return err
		}

		newLanguageList = append(newLanguageList, languageID)
	}

	if err := SyncProjectLanguages(project.ID, oldLanguageList, newLanguageList); err != nil {
		err = logger.AppendError(err)
		return err
	}

	return nil
}

func (project *Project) Get() error {
	logger := logger.Logger{Caller: "Project::Get model"}

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
		err = logger.AppendError(err)
		return err
	}

	if project.Tags, err = GetProjectTags(project.ID); err != nil {
		err = logger.AppendError(err)
		return err
	}

	if project.Languages, err = GetProjectLanguages(project.ID); err != nil {
		err = logger.AppendError(err)
		return err
	}

	return nil
}

func (project *Project) Delete() error {
	logger := logger.Logger{Caller: "Project::Delete model"}

	query := ` 
		DELETE FROM projects
		WHERE id = $1;
	`

	_, err := database.Exec(query, project.ID)
	err = logger.AppendError(err)
	return err
}

// project search related

type SortBy string

const (
	Date     SortBy = "p.created_at" //coupled with query
	Bookmark SortBy = "bookmark_count"
)

func (sortBy SortBy) value() string {
	switch sortBy {
	case Date:
		return "p.created_at"
	case Bookmark:
		return "bookmark_count"
	default:
		return "p.created_at"
	}
}

type SortDirection string

const (
	Ascending  SortDirection = "ASC"
	Descending SortDirection = "DESC"
)

func (sortDirection SortDirection) value() string {
	switch sortDirection {
	case Ascending:
		return "ASC"
	case Descending:
		return "DESC"
	default:
		return "DESC"
	}
}

type Tab string

const (
	Explore     Tab = "explore"
	MyBookmarks Tab = "my_bookmarks"
	MyProjects  Tab = "my_projects"
)

func (tabType Tab) value(userID int64) string {
	switch tabType {
	case Explore:
		return ""
	case MyBookmarks:
		return fmt.Sprintf(`AND 
			p.id IN ( SELECT project_id
								FROM project_bookmarks
								WHERE user_id = %d)`, userID)
	case MyProjects:
		return fmt.Sprintf(`AND 
				p.id IN ( SELECT id
									FROM projects
									WHERE owner_id = %d)`, userID)
	default:
		return ""
	}
}

type ProjectSearchQuery struct {
	UserID         int64
	Title          string
	TagIDs         []int64
	LanguageIDs    []int64
	OrganizationID int64
	SortBy         SortBy
	SortDirection  SortDirection
	Tab            Tab
}

func (searchQuery ProjectSearchQuery) FindProjectsWithFullTextSearch() ([]Project, error) {
	logger := logger.Logger{Caller: "ProjectSearchQuery::FindProjectsWithFullTextSearch model"}

	query := fmt.Sprintf(`
		SELECT 
			p.id,
			p.github_url,
			p.title,
			p.description,
			p.owner_id,
			p.created_at,
			COUNT(pb.user_id) AS bookmark_count
		FROM projects AS p
		LEFT JOIN project_bookmarks AS pb
		ON p.id = pb.project_id
		WHERE 
			(	$1 = '' OR
				LOWER(p.title) LIKE LOWER($1) ) AND
			( $2 = 0 OR
			  p.id IN ( SELECT project_id
							 	  FROM project_tags
					 			  WHERE tag_id IN %s ) ) AND
			( $3 = 0 OR
				p.id IN ( SELECT project_id
								  FROM project_languages
								  WHERE language_id IN %s ) ) AND
			(	$4 = 0 OR
				p.owner_id IN ( SELECT id
											  FROM users
											  WHERE organization_id = $4 ) )
			%s									
			GROUP BY p.id
			ORDER BY %s %s;`,
		ArraytoStringRoundBrackets(searchQuery.TagIDs),
		ArraytoStringRoundBrackets(searchQuery.LanguageIDs),
		searchQuery.Tab.value(searchQuery.UserID),
		searchQuery.SortBy.value(),
		searchQuery.SortDirection.value(),
	)

	args := []interface{}{
		"%" + searchQuery.Title + "%",
		len(searchQuery.TagIDs),
		len(searchQuery.LanguageIDs),
		searchQuery.OrganizationID,
	}

	rows, err := database.Query(query, args...)
	if err != nil {
		err = logger.AppendError(err)
		return nil, err
	}

	defer rows.Close()

	projects := make([]Project, 0)

	for rows.Next() {
		var project Project

		err := rows.Scan(
			&project.ID,
			&project.GithubURL,
			&project.Title,
			&project.Description,
			&project.OwnerID,
			&project.CreatedAt,
			&project.BookmarkCount,
		)

		if err != nil {
			err = logger.AppendError(err)
			return nil, err
		}

		if project.Tags, err = GetProjectTags(project.ID); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}

		if project.Languages, err = GetProjectLanguages(project.ID); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}

		projectBookmark := ProjectBookmark{UserID: searchQuery.UserID, ProjectID: project.ID}
		isBookmarked, err := projectBookmark.GetUserBookmarkStatus()
		if err != nil {
			err = logger.AppendError(err)
			return nil, err
		}
		project.UserBookmarkStatus = isBookmarked

		projects = append(projects, project)

	}
	return projects, nil
}

func GetProjectOwnerID(projectID int64) (int64, error) {
	logger := logger.Logger{Caller: "GetProjectOwnerID model"}

	query := `
		SELECT owner_id
		FROM projects
		WHERE id = $1;
	`
	var ownerID int64
	err := database.QueryRow(query, projectID).Scan(&ownerID)
	if err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return ownerID, nil
}

func IsProjectExists(projectID int64) (bool, error) {
	logger := logger.Logger{Caller: "IsProjectExists model"}

	query := `
		SELECT EXISTS (
			SELECT 1	
			FROM projects
			WHERE id = $1
		);
	`
	var exists bool
	if err := database.QueryRow(query, projectID).Scan(&exists); err != nil {
		err = logger.AppendError(err)
		return false, err
	}

	return exists, nil
}

func IsProjectExistsByURL(githubURL string) (bool, error) {
	logger := logger.Logger{Caller: "IsProjectExistsByURL model"}

	query := `
		SELECT EXISTS (
			SELECT 1	
			FROM projects
			WHERE github_url = $1
		);
	`
	var exists bool
	if err := database.QueryRow(query, githubURL).Scan(&exists); err != nil {
		err = logger.AppendError(err)
		return false, err
	}

	return exists, nil
}
