package middlewares

import (
	"net/http"
	"strconv"

	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/internal/logger"
	"github.com/priyam-trambadia/project-wall/internal/models"
)

func ValidatePathValueProjectID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.Logger{Caller: "ValidatePathValueProjectID middleware"}

		projectID, err := strconv.ParseInt(r.PathValue("project_id"), 10, 64)
		if err == strconv.ErrSyntax {
			utils.RenderInvalidProjectIDErr(w)
			return
		} else if err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}

		exists, err := models.IsProjectExists(projectID)
		if err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}

		if !exists {
			utils.RenderInvalidProjectIDErr(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func AuthorizeProjectAction(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.Logger{Caller: "AuthorizeProjectAction middleware"}

		ctx := r.Context()
		userID, _ := ctx.Value("user_id").(int64)

		projectID, _ := strconv.ParseInt(r.PathValue("project_id"), 10, 64)
		ownerID, err := models.GetProjectOwnerID(projectID)
		if err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}

		if ownerID != userID {
			utils.RenderForbiddenAccessErr(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}
