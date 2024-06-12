package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/priyam-trambadia/project-wall/web/templates"
)

func Root(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println(ctx.Value("user_id").(int64))
	templates.Home().Render(context.Background(), w)
}
