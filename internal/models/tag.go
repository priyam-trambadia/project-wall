package models

import (
	"database/sql"

	"github.com/priyam-trambadia/project-wall/internal/logger"
)

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (tag Tag) GetID() int64 {
	return tag.ID
}

func (tag *Tag) Insert() error {
	logger := logger.Logger{Caller: "Tag::Insert model"}

	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id;	
 	`
	err := database.QueryRow(query, tag.Name).Scan(&tag.ID)
	err = logger.AppendError(err)
	return err
}

func (tag *Tag) Get() error {
	logger := logger.Logger{Caller: "Tag::Get model"}

	query := `
		SELECT name
		FROM tags
		WHERE id = $1;	
 `
	err := database.QueryRow(query, tag.ID).Scan(&tag.Name)
	err = logger.AppendError(err)
	return err
}

// utilites

func GetTagID(tagName string) (int64, error) {
	logger := logger.Logger{Caller: "GetTagID model"}

	query := `
		SELECT id
		FROM tags
		WHERE name = $1;	
 `
	var tagID int64
	if err := database.QueryRow(query, tagName).Scan(&tagID); err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return tagID, nil
}

func GetOrCreateTagID(tagName string) (int64, error) {
	logger := logger.Logger{Caller: "GetOrCreateTagID model"}

	tagID, err := GetTagID(tagName)
	if err == sql.ErrNoRows {
		tag := Tag{Name: tagName}
		if err := tag.Insert(); err != nil {
			err = logger.AppendError(err)
			return 0, err
		}
		return tag.ID, nil

	} else if err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return tagID, nil
}

func FindTagsWithFullTextSearch(tagName string) ([]Tag, error) {
	logger := logger.Logger{Caller: "FindTagsWithFullTextSearch model"}

	query := `
		SELECT t.id, t.name
		FROM tags AS t
		JOIN project_tags AS pt 
		ON t.id = pt.tag_id
		WHERE LOWER(t.name) LIKE LOWER($1)
		GROUP BY t.id
		ORDER BY COUNT(pt.project_id)
		LIMIT 10;
	`

	tagNameWithWildCard := "%" + tagName + "%"

	rows, err := database.Query(query, tagNameWithWildCard)
	if err != nil {
		err = logger.AppendError(err)
		return nil, err
	}

	defer rows.Close()

	tags := make([]Tag, 0)
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func GetTagNames(tagIDs []int64) ([]string, error) {
	logger := logger.Logger{Caller: "GetTagNames model"}

	names := make([]string, 0)

	for _, tagID := range tagIDs {
		tag := Tag{ID: tagID}

		if err := tag.Get(); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}

		names = append(names, tag.Name)
	}

	return names, nil
}

// project_tag and tag JOIN

func GetProjectTags(projectID int64) ([]Tag, error) {
	logger := logger.Logger{Caller: "GetProjectTags model"}

	query := `
		SELECT t.id, t.name
		FROM tags AS t
		JOIN project_tags AS pt
		ON t.id = pt.tag_id
		WHERE pt.project_id = $1;
	`
	rows, err := database.Query(query, projectID)
	if err != nil {
		err = logger.AppendError(err)
		return nil, err
	}

	defer rows.Close()

	tags := make([]Tag, 0)
	for rows.Next() {
		var tag Tag

		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			err = logger.AppendError(err)
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
