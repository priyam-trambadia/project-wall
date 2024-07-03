package models

import "sort"

type ProjectTag struct {
	ProjectID int64
	TagID     int64
}

func (projectTag *ProjectTag) Insert() error {
	query := ` 
		INSERT INTO project_tags (project_id, tag_id)
		VALUES ($1, $2);
	`
	_, err := database.Exec(query, projectTag.ProjectID, projectTag.TagID)
	return err
}

func (projectTag *ProjectTag) Delete() error {
	query := ` 
		DELETE FROM project_tags
		WHERE project_id = $1 AND tag_id = $2;
	`
	_, err := database.Exec(query, projectTag.ProjectID, projectTag.TagID)
	return err
}

// utilites

func SyncProjectTags(projectID int64, oldTagList, newTagList []int64) error {
	sort.Slice(oldTagList, func(i, j int) bool {
		return oldTagList[i] < oldTagList[j]
	})

	sort.Slice(newTagList, func(i, j int) bool {
		return newTagList[i] < newTagList[j]
	})

	insertList := make([]int64, 0)
	deleteList := make([]int64, 0)
	oldIndex := 0
	newIndex := 0

	for oldIndex < len(oldTagList) && newIndex < len(newTagList) {
		if oldTagList[oldIndex] == newTagList[newIndex] {
			oldIndex += 1
			newIndex += 1
		} else if oldTagList[oldIndex] < newTagList[newIndex] {
			deleteList = append(deleteList, oldTagList[oldIndex])
			oldIndex += 1
		} else {
			insertList = append(insertList, newTagList[newIndex])
			newIndex += 1
		}
	}

	for oldIndex < len(oldTagList) {
		deleteList = append(deleteList, oldTagList[oldIndex])
		oldIndex += 1
	}

	for newIndex < len(insertList) {
		insertList = append(insertList, newTagList[newIndex])
		newIndex += 1
	}

	for _, tagID := range deleteList {
		deleteProjectTag := ProjectTag{ProjectID: projectID, TagID: tagID}
		if err := deleteProjectTag.Delete(); err != nil {
			return err
		}
	}

	for _, tagID := range insertList {
		insertProjectTag := ProjectTag{ProjectID: projectID, TagID: tagID}
		if err := insertProjectTag.Insert(); err != nil {
			return err
		}
	}

	return nil
}

func GetProjectTagIDs(projectID int64) ([]int64, error) {

	query := `
		SELECT tag_id
		FROM project_tags
		WHERE project_id = $1;
	`
	if rows, err := database.Query(query, projectID); err != nil {
		return nil, err
	} else {
		defer rows.Close()

		tags := make([]int64, 0)
		for rows.Next() {
			var tagID int64

			if err := rows.Scan(&tagID); err != nil {
				return nil, err
			} else {
				tags = append(tags, tagID)
			}
		}

		return tags, nil
	}
}

func GetTagProjectIDs(tagID int64) ([]int64, error) {
	query := `
		SELECT project_id
		FROM project_tags
		WHERE tag_id = $1;
	`
	if rows, err := database.Query(query, tagID); err != nil {
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
