package models

import "fmt"

type Tag struct {
	ID   int64
	Name string
}

func (tag *Tag) Insert() {
	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id;	
 	`
	database.QueryRow(query, tag.Name).Scan(&tag.ID)
}

func (tag *Tag) GetID() error {
	return fmt.Errorf("Tag Not Found")
}
