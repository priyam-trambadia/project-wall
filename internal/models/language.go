package models

import "database/sql"

type Language struct {
	ID   int64
	Name string
}

func (language *Language) Insert() error {
	query := `
		INSERT INTO languages (name)
		VALUES ($1)
		RETURNING id;	
 	`
	return database.QueryRow(query, language.Name).Scan(&language.ID)
}

func (language *Language) Get() error {
	query := `
		SELECT name
		FROM languages
		WHERE id = $1;
 `
	return database.QueryRow(query, language.ID).Scan(&language.Name)
}

// utilites

func GetLanguageID(languageName string) (int64, error) {
	query := `
		SELECT id
		FROM languages
		WHERE name = $1;
 	`
	var languageID int64
	if err := database.QueryRow(query, languageName).Scan(&languageID); err != nil {
		return 0, err
	} else {
		return languageID, nil
	}
}

func GetOrCreateLanguageID(languageName string) (int64, error) {
	if languageID, err := GetLanguageID(languageName); err == sql.ErrNoRows {
		language := Language{Name: languageName}
		if err := language.Insert(); err != nil {
			return 0, err
		} else {
			return language.ID, nil
		}
	} else if err != nil {
		return 0, err
	} else {
		return languageID, nil
	}
}

func FindLanguagesWithFullTextSearch(languageName string) ([]Language, error) {
	query := `
		SELECT id, name
		FROM languages
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '');
	`

	if rows, err := database.Query(query, languageName); err != nil {
		return nil, err
	} else {
		languages := make([]Language, 0)
		for rows.Next() {
			var language Language
			if err := rows.Scan(&language.ID, &language.Name); err != nil {
				return nil, err
			} else {
				languages = append(languages, language)
			}
		}

		return languages, nil
	}
}

func GetLanguageNames(languageIDs []int64) ([]string, error) {
	names := make([]string, 0)

	for _, languageID := range languageIDs {
		language := Language{ID: languageID}

		if err := language.Get(); err != nil {
			return nil, err
		} else {
			names = append(names, language.Name)
		}
	}

	return names, nil
}
