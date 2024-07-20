package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/internal/models"
	"github.com/priyam-trambadia/project-wall/web/templates"
)

func Root(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isUserLogin, _ := ctx.Value("is_user_logged_in").(bool)
	userID, _ := ctx.Value("user_id").(int64)

	projectSearchQuery := models.ProjectSearchQuery{UserID: userID}
	projects, err := projectSearchQuery.FindProjectsWithFullTextSearch()
	if err != nil {
		utils.RenderInternalServerErr(w)
		log.Fatalln("[-] Error in Root handler.\n", err)
	}

	tags, err := models.FindTagsWithFullTextSearch("")
	languages, err := models.FindLanguagesWithFullTextSearch("")

	templates.HomePage(isUserLogin, userID, projects, tags, languages).Render(context.Background(), w)
}
