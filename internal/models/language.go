package models

import (
	"database/sql"

	"github.com/priyam-trambadia/project-wall/internal/logger"
)

type Language struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (language Language) GetID() int64 {
	return language.ID
}

func (language *Language) Insert() error {
	logger := logger.Logger{Caller: "Language::Insert model"}
	query := `
		INSERT INTO languages (name)
		VALUES ($1)
		RETURNING id;	
 	`
	err := database.QueryRow(query, language.Name).Scan(&language.ID)
	if err != nil {
		err = logger.AppendError(err)
	}
	return err
}

func (language *Language) Get() error {
	logger := logger.Logger{Caller: "Language::Get model"}
	query := `
		SELECT name
		FROM languages
		WHERE id = $1;
 `
	err := database.QueryRow(query, language.ID).Scan(&language.Name)
	if err != nil {
		err = logger.AppendError(err)
	}
	return err
}

// utilites
func GetLanguageID(languageName string) (int64, error) {
	logger := logger.Logger{Caller: "GetLanguageID model"}
	query := `
		SELECT id
		FROM languages
		WHERE name = $1;
 	`
	var languageID int64
	if err := database.QueryRow(query, languageName).Scan(&languageID); err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return languageID, nil
}

func GetOrCreateLanguageID(languageName string) (int64, error) {
	logger := logger.Logger{Caller: "GetOrCreateLanguageID model"}

	languageID, err := GetLanguageID(languageName)
	if err == sql.ErrNoRows {

		language := Language{Name: languageName}
		if err := language.Insert(); err != nil {
			err = logger.AppendError(err)
			return 0, err
		}
		return language.ID, nil

	} else if err != nil {

		err = logger.AppendError(err)
		return 0, err
	}

	return languageID, nil

}

func FindLanguagesWithFullTextSearch(languageName string) ([]Language, error) {
	logger := logger.Logger{Caller: "FindLanguagesWithFullTextSearch model"}

	query := `
		SELECT l.id, l.name
		FROM languages AS l
		JOIN project_languages AS pl
		ON l.id = pl.language_id
		WHERE LOWER(l.name) LIKE LOWER($1)
		GROUP BY l.id
		ORDER BY COUNT(pl.project_id)
		LIMIT 10;
	`

	languageNameWithWildCard := "%" + languageName + "%"

	rows, err := database.Query(query, languageNameWithWildCard)
	if err != nil {
		err = logger.AppendError(err)
		return nil, err
	}
	defer rows.Close()

	languages := make([]Language, 0)
	for rows.Next() {
		var language Language
		if err := rows.Scan(&language.ID, &language.Name); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}
		languages = append(languages, language)
	}

	return languages, nil
}

func GetLanguageNames(languageIDs []int64) ([]string, error) {
	logger := logger.Logger{Caller: "GetLanguageNames model"}

	names := make([]string, 0)

	for _, languageID := range languageIDs {
		language := Language{ID: languageID}

		if err := language.Get(); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}
		names = append(names, language.Name)
	}

	return names, nil
}

// project_language and language JOIN

func GetProjectLanguages(projectID int64) ([]Language, error) {
	logger := logger.Logger{Caller: "GetProjectLanguages model"}

	query := `
		SELECT l.id, l.name
		FROM languages AS l
		JOIN project_languages AS pl
		ON l.id = pl.language_id
		WHERE pl.project_id = $1;
	`
	rows, err := database.Query(query, projectID)
	if err != nil {
		err = logger.AppendError(err)
		return nil, err
	}

	defer rows.Close()

	languages := make([]Language, 0)
	for rows.Next() {
		var language Language

		if err := rows.Scan(&language.ID, &language.Name); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}
		languages = append(languages, language)

	}

	return languages, nil
}
