package controllers

import (
	"context"
	"net/http"

	"github.com/priyam-trambadia/project-wall/src/ui/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	pages.Home().Render(context.Background(), w)
}
