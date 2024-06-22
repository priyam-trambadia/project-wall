package models

import "time"

type Project struct {
	ID            int64
	GithubLink    string
	Title         string
	Description   string
	OwnerID       int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastSyncedAt  time.Time
	Tags          []string
	Languages     []string
	BookmarkCount int64
}

func (project *Project) Insert() {
	query := `
		INSERT INTO projects (github_link, title, description, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	args := []interface{}{
		project.GithubLink,
		project.Title,
		project.Description,
		project.OwnerID,
	}

	database.QueryRow(query, args...).Scan(&project.ID)

	//adding tags
	for _, tagName := range project.Tags {
		tag := Tag{Name: tagName}

		err := tag.GetID()
		if err != nil {
			tag.Insert()
		}

		projectTag := ProjectTag{ProjectID: project.ID, TagID: tag.ID}
		projectTag.Insert()
	}

	//adding languages
	for _, languageName := range project.Languages {
		language := Language{Name: languageName}

		err := language.GetID()
		if err != nil {
			language.Insert()
		}

		projectLanguage := ProjectLanguage{ProjectID: project.ID, LanguageID: language.ID}
		projectLanguage.Insert()
	}
}

func (project *Project) Update() {
	query := `
		UPDATE projects
		SET title = $1, description = $2, updated_at = $3,  
		WHERE id = $4;
	`
	args := []interface{}{
		project.Title,
		project.Description,
		time.Now(),
		project.ID,
	}

	database.Exec(query, args...)

	// updating tags
	projectTag := ProjectTag{ProjectID: project.ID}
	oldTagList := projectTag.GetTags()

	newTagList := make([]int64, 0)
	for _, tagName := range project.Tags {
		tag := Tag{Name: tagName}

		err := tag.GetID()
		if err != nil {
			tag.Insert()
		}

		newTagList = append(newTagList, tag.ID)
	}

	projectTag.UpdateProjectTags(oldTagList, newTagList)

	// updating languages
	projectLanguage := ProjectLanguage{ProjectID: project.ID}
	oldLanguageList := projectLanguage.GetLanguages()

	newLanguageList := make([]int64, 0)
	for _, languageName := range project.Languages {
		language := Language{Name: languageName}

		err := language.GetID()
		if err != nil {
			language.Insert()
		}

		newLanguageList = append(newLanguageList, language.ID)
	}

	projectLanguage.UpdateProjectLanguages(oldLanguageList, newLanguageList)
}

func (project *Project) Delete() {
	query := ` 
		DELETE FROM projects
		WHERE project_id = $1;
	`

	database.Exec(query, project.ID)
}

func GetAllProjects() []Project {
	projects := make([]Project, 0)

	return projects
}
