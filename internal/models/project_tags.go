package models

import "sort"

type ProjectTag struct {
	ProjectID int64
	TagID     int64
}

func (projectTag *ProjectTag) Insert() {
	query := ` 
		INSERT INTO project_tags (project_id, tag_id)
		VALUES ($1, $2);
	`
	database.QueryRow(query, projectTag.ProjectID, projectTag.TagID)
}

func (projectTag *ProjectTag) Delete() {
	query := ` 
		DELETE FROM project_tags
		WHERE project_id = $1 AND tag_id = $2;
	`

	database.QueryRow(query, projectTag.ProjectID, projectTag.TagID)
}

func (projectTag *ProjectTag) GetTags() []int64 {
	tags := make([]int64, 0)

	query := `
		SELECT tag_id
		FROM project_tags
		WHERE project_id = $1;
	`

	rows, _ := database.Query(query, projectTag.ProjectID)
	defer rows.Close()

	for rows.Next() {
		var tagID int64

		err := rows.Scan(&tagID)
		if err != nil {
			break
		}

		tags = append(tags, tagID)
	}

	return tags
}

func (projectTag *ProjectTag) GetProjects() []int64 {
	projects := make([]int64, 0)

	query := `
		SELECT project_id
		FROM project_tags
		WHERE tag_id = $1;
	`
	rows, _ := database.Query(query, projectTag.TagID)
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

func (projectTag *ProjectTag) UpdateProjectTags(oldTagList []int64, newTagList []int64) {
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
		deleteProjectTag := ProjectTag{ProjectID: projectTag.ProjectID, TagID: tagID}
		deleteProjectTag.Delete()

	}

	for _, tagID := range insertList {
		insertProjectTag := ProjectTag{ProjectID: projectTag.ProjectID, TagID: tagID}
		insertProjectTag.Insert()
	}
}
