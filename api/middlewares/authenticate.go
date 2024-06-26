package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/priyam-trambadia/project-wall/api/handlers"
	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ctx            context.Context
			isUserLoggedIn bool
			userID         int64
			ok             bool
		)

		accessTokenCookie, err := r.Cookie("access_token")
		if err == nil {
			userID, ok = jwt.ValidateAccessToken(accessTokenCookie.Value)
		}

		if err != nil || !ok {

			refreshTokenCookie, err2 := r.Cookie("refresh_token")
			if err2 == nil {
				userID, ok = jwt.ValidateRefreshToken(refreshTokenCookie.Value)
			}

			if err2 != nil || !ok {
				isUserLoggedIn = false

			} else {
				isUserLoggedIn = true
				accessToken := jwt.GenerateAccessToken(userID)
				utils.SetTokenCookie(w, accessToken, refreshTokenCookie.Value)
			}

		} else {
			isUserLoggedIn = true
		}

		ctx = context.WithValue(r.Context(), "is_user_logged_in", isUserLoggedIn)
		ctx = context.WithValue(ctx, "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthenticationRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isUserLoggedIn, ok := ctx.Value("is_user_logged_in").(bool)

		if !ok {
			log.Println("[-] Error in AuthenticationRequired type conversion.")
		}

		if !ok || !isUserLoggedIn {
			handlers.UserLogin(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}
}
