package middlewares

import (
	"net/http"
	"strconv"

	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/internal/logger"
	"github.com/priyam-trambadia/project-wall/internal/models"
)

func ValidatePathValueUserID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.Logger{Caller: "ValidatePathValueUserID middleware"}

		userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
		if err == strconv.ErrSyntax {
			utils.RenderInvalidUserIDErr(w)
			return
		} else if err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}

		exists, err := models.IsUserExists(userID)
		if err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}

		if !exists {
			utils.RenderInvalidUserIDErr(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func AuthorizeUserProfileAction(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIDContext, _ := ctx.Value("user_id").(int64)
		userIDPathValue, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)

		if userIDContext != userIDPathValue {
			utils.RenderForbiddenAccessErr(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}
