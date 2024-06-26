package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/priyam-trambadia/project-wall/internal/models"
	"github.com/priyam-trambadia/project-wall/web/templates"
)

func ProjectAdd(w http.ResponseWriter, r *http.Request) {
	projectJSONString := r.FormValue("project_details_json")

	var project models.Project

	// Unmarshal the JSON string into the Person struct
	err := json.Unmarshal([]byte(projectJSONString), &project)
	if err != nil {
		// fmt.Println("Error unmarshaling JSON:", err)
		// return
	}

	languageJSON, _ := json.Marshal(project.Languages)
	templates.AddProjectPage(project, string(languageJSON)).Render(context.Background(), w)

}

func ProjectAddPOST(w http.ResponseWriter, r *http.Request) {
	var project models.Project

	ctx := r.Context()

	project.OwnerID = ctx.Value("user_id").(int64)
	project.GithubLink = r.FormValue("github_link")
	project.Title = r.FormValue("title")
	project.Description = r.FormValue("description")
	json.Unmarshal([]byte(r.FormValue("languages")), &project.Languages)
	json.Unmarshal([]byte(r.FormValue("tags")), &project.Tags)

	fmt.Println(project.OwnerID)
	project.Insert()

}

func ProjectTagSearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("tags"))

	w.Write([]byte(
		`<span class="clickable-span"> hi </span>` +
			`<span class="clickable-span"> hi2 </span>` +
			`<span class="clickable-span"> hi3 </span>`,
	))
}

func ProjectLanguageSearch(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(
		`<span> Lang </span>`,
	))
}
