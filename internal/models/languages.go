package models

type Language struct {
	ID   int64
	Name string
}

func (language *Language) Insert() {
	query := `
		INSERT INTO languages (name)
		VALUES ($1)
		RETURNING id;	
 	`
	database.QueryRow(query, language.Name).Scan(&language.ID)
}

func (language *Language) GetID() error {

	return nil
}
