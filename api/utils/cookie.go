package utils

import (
	"net/http"
	"time"

	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
)

func SetTokenCookie(w http.ResponseWriter, accessToken string, refreshToken string) {
	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Expires:  time.Now().Add(jwt.AccessTokenExpiry),
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Expires:  time.Now().Add(jwt.RefreshTokenExpiry),
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

}

func DeleteTokenCookie(w http.ResponseWriter) {
	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
}
