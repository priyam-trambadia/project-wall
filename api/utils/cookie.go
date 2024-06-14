package utils

import "net/http"

func SetTokenCookie(w http.ResponseWriter, accessToken string, refreshToken string) {
	accessTokenCookie := &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
		Path:  "/",
	}

	refreshTokenCookie := &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		Path:  "/",
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

}

func DeleteTokenCookie(w http.ResponseWriter) {
	accessTokenCookie := &http.Cookie{
		Name:   "access_token",
		MaxAge: -1,
		Path:   "/",
	}

	refreshTokenCookie := &http.Cookie{
		Name:   "refresh_token",
		MaxAge: -1,
		Path:   "/",
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
}
