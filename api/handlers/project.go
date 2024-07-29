package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/internal/logger"
	"github.com/priyam-trambadia/project-wall/internal/models"
	"github.com/priyam-trambadia/project-wall/web/templates"
)

func ProjectCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	isUserLogin, _ := ctx.Value("is_user_login").(bool)
	userID, _ := ctx.Value("user_id").(int64)

	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	projectJSONString := r.FormValue("project_detail_json")
	var project models.Project

	err := json.Unmarshal([]byte(projectJSONString), &project)
	if err != nil {
		utils.RenderInvalidJSONErr(w)
		return
	}

	templates.AddProjectPage(isUserLogin, userID, project).Render(context.Background(), w)
}

func ProjectCreatePOST(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "ProjectCreate handler"}

	var project models.Project

	ctx := r.Context()
	project.OwnerID = ctx.Value("user_id").(int64)

	project.GithubURL = r.FormValue("github_url")
	project.Title = r.FormValue("title")
	project.Description = r.FormValue("description")

	if err := json.Unmarshal([]byte(r.FormValue("languages")), &project.Languages); err != nil {
		utils.RenderInvalidJSONErr(w)
		logger.Println(err)
		return
	}

	if err := json.Unmarshal([]byte(r.FormValue("tags")), &project.Tags); err != nil {
		utils.RenderInvalidJSONErr(w)
		logger.Println(err)
		return
	}

	exists, err := models.IsProjectExistsByURL(project.GithubURL)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	if exists {
		utils.SetPopupCookie(w, "This project already uploaded.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := project.Insert(); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	utils.SetPopupCookie(w, "Project added successfully.")
	http.Redirect(w, r, "/", http.StatusFound)
}

func ProjectUpdate(w http.ResponseWriter, r *http.Request) {

}

func ProjectDelete(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "ProjectDelete handler"}
	projectID, _ := strconv.ParseInt(r.PathValue("project_id"), 10, 64)

	project := models.Project{ID: projectID}
	if err := project.Delete(); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}
}

func ProjectToggleBookmark(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "ProjectToggleBookmark handler"}

	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(int64)
	projectID, _ := strconv.ParseInt(r.PathValue("project_id"), 10, 64)

	projectBookmark := models.ProjectBookmark{ProjectID: projectID, UserID: userID}
	isBookmarked, err := projectBookmark.GetUserBookmarkStatus()
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	if isBookmarked {
		if err := projectBookmark.Delete(); err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}
	} else {
		if err = projectBookmark.Insert(); err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}
	}
}

func ProjectSearch(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "ProjectSearch handler"}

	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(int64)

	var searchQuery models.ProjectSearchQuery

	searchQuery.UserID, _ = strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	searchQuery.Title = r.FormValue("project-search")

	switch r.FormValue("sort-by") {
	case "date":
		searchQuery.SortBy = models.Date
	case "bookmark":
		searchQuery.SortBy = models.Bookmark
	}

	switch r.FormValue("sort-direction") {
	case "asc":
		searchQuery.SortDirection = models.Ascending
	case "desc":
		searchQuery.SortDirection = models.Descending
	}

	switch r.FormValue("tab") {
	case "explore":
		searchQuery.Tab = models.Explore
	case "my_bookmarks":
		searchQuery.Tab = models.MyBookmarks
	case "my_projects":
		searchQuery.Tab = models.MyProjects
	}

	if r.FormValue("organization-only") == "true" {
		searchQuery.OrganizationID, _ = models.GetUserOrganizationID(userID)
	}

	languages := make([]models.Language, 0)
	if err := json.Unmarshal([]byte(r.FormValue("languages")), &languages); err != nil {
		utils.RenderInvalidJSONErr(w)
		return
	}
	searchQuery.LanguageIDs = utils.ExtractIDsMakeArray(languages)

	tags := make([]models.Tag, 0)
	if err := json.Unmarshal([]byte(r.FormValue("tags")), &tags); err != nil {
		utils.RenderInvalidJSONErr(w)
		return
	}
	searchQuery.TagIDs = utils.ExtractIDsMakeArray(tags)

	projects, err := searchQuery.FindProjectsWithFullTextSearch()
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	for _, project := range projects {
		templates.ProjectCard(userID, project).Render(context.Background(), w)
	}
}

func ProjectTagSearch(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "ProjectTagSearch handler"}
	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	tagName := r.FormValue("tag-search")

	tags, err := models.FindTagsWithFullTextSearch(tagName)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	if include, err := strconv.ParseBool(r.FormValue("include_same_tag_if_not_exists")); err != nil {
		utils.RenderFormParsingErr(w)
		return
	} else if include && tagName != "" {
		isFound := false
		for _, tag := range tags {
			if tag.Name == tagName {
				isFound = true
				break
			}
		}
		if !isFound {
			tags = append(tags, models.Tag{ID: 0, Name: tagName})
		}
	}

	for _, tag := range tags {
		templates.ClickableSpanLi(tag.ID, tag.Name).Render(context.Background(), w)
	}
}

func ProjectLanguageSearch(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "ProjectLanguageSearch handler"}
	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	languagesName := r.FormValue("language-search")

	languages, err := models.FindLanguagesWithFullTextSearch(languagesName)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	for _, language := range languages {
		templates.ClickableSpanLi(language.ID, language.Name).Render(context.Background(), w)
	}
}
