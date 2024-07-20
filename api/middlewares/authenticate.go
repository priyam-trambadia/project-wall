package middlewares

import (
	"context"
	"net/http"

	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
	"github.com/priyam-trambadia/project-wall/internal/logger"
	"github.com/priyam-trambadia/project-wall/internal/models"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger := logger.Logger{Caller: "Authenticate middleware"}

		var (
			userID         int64 = 0
			isUserLoggedIn bool  = false
			ok             bool  = false
			err            error
		)

		accessTokenCookie, _ := r.Cookie("access_token")

		if accessTokenCookie != nil {
			userID, ok, err = jwt.ValidateAccessToken(accessTokenCookie.Value)
			if err != nil {
				utils.RenderInternalServerErr(w)
				logger.Fatalln(err)
			}
		}

		if ok {
			isUserLoggedIn = true
		} else {

			refreshTokenCookie, _ := r.Cookie("refresh_token")
			if refreshTokenCookie != nil {
				userID, ok, err = jwt.ValidateRefreshToken(refreshTokenCookie.Value)
				if err != nil {
					utils.RenderInternalServerErr(w)
					logger.Fatalln(err)
				}
			}

			if ok {
				isUserLoggedIn = true
				var accessToken string
				accessToken, err := jwt.GenerateAccessToken(userID)
				if err != nil {
					utils.RenderInternalServerErr(w)
					logger.Fatalln(err)
				}
				utils.SetTokenCookie(w, accessToken, refreshTokenCookie.Value)
			}
		}

		var ctx context.Context
		ctx = context.WithValue(r.Context(), "is_user_logged_in", isUserLoggedIn)
		ctx = context.WithValue(ctx, "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserAuthenticationRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.Logger{Caller: "UserAuthenticationRequired middleware"}

		ctx := r.Context()
		isUserLoggedIn, _ := ctx.Value("is_user_logged_in").(bool)

		if !isUserLoggedIn {
			utils.SetPopupCookie(w, "You need to be logged in to continue. Please log in or create an account.")
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		userID, _ := ctx.Value("user_id").(int64)
		exists, err := models.IsUserExists(userID)
		if err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}

		if !exists {
			utils.RenderSessionTemperedErr(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}
