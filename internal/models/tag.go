package models

import "database/sql"

type Tag struct {
	ID   int64
	Name string
}

func (tag *Tag) Insert() error {
	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id;	
 	`
	return database.QueryRow(query, tag.Name).Scan(&tag.ID)
}

func (tag *Tag) Get() error {
	query := `
		SELECT name
		FROM tags
		WHERE id = $1;	
 `
	return database.QueryRow(query, tag.ID).Scan(&tag.Name)
}

// utilites

func GetTagID(tagName string) (int64, error) {
	query := `
		SELECT id
		FROM tags
		WHERE name = $1;	
 `
	var tagID int64
	if err := database.QueryRow(query, tagName).Scan(&tagID); err != nil {
		return 0, err
	} else {
		return tagID, nil
	}
}

func GetOrCreateTagID(tagName string) (int64, error) {
	if tagID, err := GetTagID(tagName); err == sql.ErrNoRows {
		tag := Tag{Name: tagName}
		if err := tag.Insert(); err != nil {
			return 0, err
		} else {
			return tag.ID, nil
		}
	} else if err != nil {
		return 0, err
	} else {
		return tagID, nil
	}
}

func FindTagsWithFullTextSearch(tagName string) ([]Tag, error) {
	query := `
		SELECT id, name
		FROM tags
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '');
	`

	if rows, err := database.Query(query, tagName); err != nil {
		return nil, err
	} else {
		tags := make([]Tag, 0)
		for rows.Next() {
			var tag Tag
			if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
				return nil, err
			} else {
				tags = append(tags, tag)
			}
		}

		return tags, nil
	}
}

func GetTagNames(tagIDs []int64) ([]string, error) {
	names := make([]string, 0)

	for _, tagID := range tagIDs {
		tag := Tag{ID: tagID}

		if err := tag.Get(); err != nil {
			return nil, err
		} else {
			names = append(names, tag.Name)
		}
	}

	return names, nil
}
