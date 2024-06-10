package handlers

import (
	"context"
	"net/http"

	"github.com/priyam-trambadia/project-wall/internal/models"
	"github.com/priyam-trambadia/project-wall/web/templates"
)

func Root(w http.ResponseWriter, r *http.Request) {
	var user models.User
	user.Insert()

	templates.Home().Render(context.Background(), w)
}
