package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
	"github.com/priyam-trambadia/project-wall/internal/models"
	"github.com/priyam-trambadia/project-wall/web/templates"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	isLogin := ctx.Value("is_user_logged_in").(bool)

	if isLogin {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		templates.LoginPage().Render(context.Background(), w)
	}
}

func UserLoginPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // Parse form data from the request body
	if err != nil {
		fmt.Println(err)
		// Handle error (e.g., log the error or return an error response)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user := models.User{Email: email, Password: password}
	ok := user.ValidateUser()

	if !ok {
		fmt.Println("not ok\n", user)
	} else {
		accessToken := jwt.GenerateAccessToken(user.ID)
		user.RefreshToken = jwt.GenerateRefreshToken(user.ID)

		user.UpdateRefreshToken()

		utils.SetTokenCookie(w, accessToken, user.RefreshToken)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	isLogin := ctx.Value("is_user_logged_in").(bool)

	if isLogin {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		templates.RegisterPage().Render(context.Background(), w)
	}
}

func UserRegisterPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // Parse form data from the request body
	if err != nil {
		// Handle error (e.g., log the error or return an error response)
		return
	}

	var user models.User
	user.Name = r.FormValue("name")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	user.Insert()

	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	utils.DeleteTokenCookie(w)
	http.Redirect(w, r, "/", http.StatusFound)

}
