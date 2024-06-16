package models

import "sort"

type ProjectLanguage struct {
	ProjectID  int64
	LanguageID int64
}

func (projectLanguage *ProjectLanguage) Insert() {
	query := ` 
		INSERT INTO project_languages (project_id, language_id)
		VALUES ($1, $2);
	`
	database.QueryRow(query, projectLanguage.ProjectID, projectLanguage.LanguageID)
}

func (projectLanguage *ProjectLanguage) Delete() {
	query := ` 
		DELETE FROM project_languages
		WHERE project_id = $1 AND language_id = $2;
	`

	database.QueryRow(query, projectLanguage.ProjectID, projectLanguage.LanguageID)
}

func (projectLanguage *ProjectLanguage) GetLanguages() []int64 {
	languages := make([]int64, 0)

	query := `
		SELECT language_id
		FROM project_languages
		WHERE project_id = $1;
	`

	rows, _ := database.Query(query, projectLanguage.ProjectID)
	defer rows.Close()

	for rows.Next() {
		var LanguageID int64

		err := rows.Scan(&LanguageID)
		if err != nil {
			break
		}

		languages = append(languages, LanguageID)
	}

	return languages
}

func (projectLanguage *ProjectLanguage) GetProjects() []int64 {
	projects := make([]int64, 0)

	query := `
		SELECT project_id
		FROM project_languages
		WHERE language_id = $1;
	`
	rows, _ := database.Query(query, projectLanguage.LanguageID)
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

func (projectLanguage *ProjectLanguage) UpdateProjectLanguages(oldLanguageList []int64, newLanguageList []int64) {
	sort.Slice(oldLanguageList, func(i, j int) bool {
		return oldLanguageList[i] < oldLanguageList[j]
	})

	sort.Slice(newLanguageList, func(i, j int) bool {
		return newLanguageList[i] < newLanguageList[j]
	})

	insertList := make([]int64, 0)
	deleteList := make([]int64, 0)
	oldIndex := 0
	newIndex := 0

	for oldIndex < len(oldLanguageList) && newIndex < len(newLanguageList) {
		if oldLanguageList[oldIndex] == newLanguageList[newIndex] {
			oldIndex += 1
			newIndex += 1
		} else if oldLanguageList[oldIndex] < newLanguageList[newIndex] {
			deleteList = append(deleteList, oldLanguageList[oldIndex])
			oldIndex += 1
		} else {
			insertList = append(insertList, newLanguageList[newIndex])
			newIndex += 1
		}
	}

	for oldIndex < len(oldLanguageList) {
		deleteList = append(deleteList, oldLanguageList[oldIndex])
		oldIndex += 1
	}

	for newIndex < len(insertList) {
		insertList = append(insertList, newLanguageList[newIndex])
		newIndex += 1
	}

	for _, LanguageID := range deleteList {
		deleteProjectLanguage := ProjectLanguage{ProjectID: projectLanguage.ProjectID, LanguageID: LanguageID}
		deleteProjectLanguage.Delete()

	}

	for _, LanguageID := range insertList {
		insertProjectLanguage := ProjectLanguage{ProjectID: projectLanguage.ProjectID, LanguageID: LanguageID}
		insertProjectLanguage.Insert()
	}
}
