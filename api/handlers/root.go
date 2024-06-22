package handlers

import (
	"context"
	"net/http"

	"github.com/priyam-trambadia/project-wall/internal/models"
	"github.com/priyam-trambadia/project-wall/web/templates"
)

func Root(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isLogin := ctx.Value("is_user_logged_in").(bool)
	projects := models.GetAllProjects()
	templates.HomePage(isLogin, projects).Render(context.Background(), w)
}
