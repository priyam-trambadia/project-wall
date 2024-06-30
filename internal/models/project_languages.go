package models

import "sort"

type ProjectLanguage struct {
	ProjectID  int64
	LanguageID int64
}

func (projectLanguage *ProjectLanguage) Insert() error {
	query := ` 
		INSERT INTO project_languages (project_id, language_id)
		VALUES ($1, $2);
	`
	_, err := database.Exec(query, projectLanguage.ProjectID, projectLanguage.LanguageID)
	return err
}

func (projectLanguage *ProjectLanguage) Delete() error {
	query := ` 
		DELETE FROM project_languages
		WHERE project_id = $1 AND language_id = $2;
	`

	_, err := database.Exec(query, projectLanguage.ProjectID, projectLanguage.LanguageID)
	return err
}

// utilities

func SyncProjectLanguages(projectID int64, oldLanguageList, newLanguageList []int64) error {
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
		deleteProjectLanguage := ProjectLanguage{ProjectID: projectID, LanguageID: LanguageID}
		if err := deleteProjectLanguage.Delete(); err != nil {
			return err
		}
	}

	for _, LanguageID := range insertList {
		insertProjectLanguage := ProjectLanguage{ProjectID: projectID, LanguageID: LanguageID}
		if err := insertProjectLanguage.Insert(); err != nil {
			return err
		}
	}

	return nil
}

func GetProjectLanguagesIDs(projectID int64) ([]int64, error) {
	query := `
		SELECT language_id
		FROM project_languages
		WHERE project_id = $1;
	`
	if rows, err := database.Query(query, projectID); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		languages := make([]int64, 0)

		for rows.Next() {
			var languageID int64

			if err := rows.Scan(&languageID); err != nil {
				return nil, err
			} else {

				languages = append(languages, languageID)
			}
		}
		return languages, nil
	}
}

func GetLanguageProjectIDs(languageID int64) ([]int64, error) {
	query := `
		SELECT project_id
		FROM project_languages
		WHERE language_id = $1;
	`
	if rows, err := database.Query(query, languageID); err != nil {
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
