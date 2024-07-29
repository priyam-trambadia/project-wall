package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/priyam-trambadia/project-wall/api/utils"
	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
	"github.com/priyam-trambadia/project-wall/internal/logger"
	"github.com/priyam-trambadia/project-wall/internal/mailer"
	"github.com/priyam-trambadia/project-wall/internal/models"
	"github.com/priyam-trambadia/project-wall/web/templates"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	templates.RegisterPage().Render(context.Background(), w)
}

func UserRegisterPOST(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserRegisterPOST handler"}

	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	var user models.User
	user.Name = r.FormValue("name")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	if ok, err := models.IsEmailExists(user.Email); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	} else if ok {
		utils.SetPopupCookie(w, "A user is already registered with this email address.")
		http.Redirect(w, r, "/user/register", http.StatusFound)
		return
	}

	if err := user.Insert(); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	token, err := jwt.GenerateUserActivationToken(user.ID)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	link := utils.GetUserActivationLink(token)
	if err := mailer.SendUserActivationMail(user.Email, link); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	utils.SetPopupCookie(w, "We sent an activation email. Please check your inbox.")
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	templates.LoginPage().Render(context.Background(), w)
}

func UserLoginPOST(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserLoginPOST handler"}

	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	userID, err := models.GetUserID(email, password)
	if err == models.ErrRecordNotFound {
		utils.SetPopupCookie(w, "The email address and/or password you specified are not correct.")
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	} else if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	user := models.User{ID: userID}
	user.Get()

	if !user.IsActivated {
		token, err := jwt.GenerateUserActivationToken(user.ID)
		if err != nil {
			utils.RenderInternalServerErr(w)
			logger.Fatalln(err)
		}

		link := utils.GetUserActivationLink(token)
		mailer.SendUserActivationMail(user.Email, link)
		utils.SetPopupCookie(
			w,
			"Your account is not yet activated. We have re-sent the activation email. Please check your inbox.",
		)
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	}

	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	if err := models.UpdateUserRefreshToken(userID, refreshToken); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	utils.SetTokenCookie(w, accessToken, refreshToken)
	http.Redirect(w, r, "/", http.StatusFound)
}

func UserActivate(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserActivate handler"}
	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	token := r.FormValue("token")
	userID, ok, err := jwt.ValidateUserActivationToken(token)

	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	if !ok {
		utils.RenderInvalidTokenErr(w)
		return
	}

	if err := models.ActivateUser(userID); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	utils.SetPopupCookie(w, "Your account activated successfully.")
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func UserForgotPassword(w http.ResponseWriter, r *http.Request) {
	templates.ForgotPasswordPage().Render(context.Background(), w)
}

func UserForgotPasswordPOST(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserForgotPasswordPOST handler"}

	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	email := r.FormValue("email")
	userID, err := models.GetUserIDByEmail(email)
	if err == models.ErrRecordNotFound {
		utils.SetPopupCookie(w, "Email not found in our database")
		http.Redirect(w, r, "/user/forgot-password", http.StatusFound)
		return
	} else if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	token, err := jwt.GenerateUserPasswordResetToken(userID)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	link := utils.GetUserPasswordResetLink(token)
	if err := mailer.SendUserPasswordResetMail(email, link); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	utils.SetPopupCookie(w, "We have sent you password reset mail. Check your inbox.")
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func UserPasswordReset(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	token := r.FormValue("token")
	templates.ResetPasswordPage(token).Render(context.Background(), w)
}

func UserPasswordResetPOST(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserPasswordResetPOST handler"}

	if err := r.ParseForm(); err != nil {
		utils.RenderFormParsingErr(w)
		return
	}

	token := r.FormValue("token")
	newPassword := r.FormValue("password")

	userID, ok, err := jwt.ValidateUserPasswordResetToken(token)
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	if !ok {
		utils.RenderInvalidTokenErr(w)
		return
	}

	if err := models.UpdateUserPassword(userID, newPassword); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	utils.SetPopupCookie(w, "Your password has been change successfully.")
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserLogout handler"}

	utils.DeleteTokenCookie(w)

	ctx := r.Context()
	userID, _ := ctx.Value("user_id").(int64)

	if err := models.UpdateUserRefreshToken(userID, jwt.ValEmptyToken); err == models.ErrRecordNotFound {
		utils.SetPopupCookie(w, "No user found with given details.")
	} else if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserAvatar(w http.ResponseWriter, r *http.Request) {
}

func UserGetProfile(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserGetProfile handler"}

	ctx := r.Context()
	isUserLogin, _ := ctx.Value("is_user_logged_in").(bool)
	userID, _ := ctx.Value("user_id").(int64)

	var user models.User
	user.ID, _ = strconv.ParseInt(r.PathValue("user_id"), 10, 64)

	if err := user.Get(); err == models.ErrRecordNotFound {
		utils.SetPopupCookie(w, "User details not found.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	var searchQuery models.ProjectSearchQuery
	searchQuery.UserID = user.ID
	searchQuery.Tab = models.MyProjects

	projects, err := searchQuery.FindProjectsWithFullTextSearch()
	if err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	templates.UserProfilePage(isUserLogin, userID, user, projects).Render(context.Background(), w)
}

func UserUpdateProfile(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserUpdateProfile handler"}

	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	user := models.User{ID: userID}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err := user.Get(); err != nil {
		logger.Fatalln(err)
	}

	if err := json.Unmarshal(body, &user); err != nil {
		utils.RenderInvalidJSONErr(w)
		return
	}

	if err := user.Update(); err != nil {
		logger.Fatalln(err)
	}

}

func UserDeleteProfile(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger{Caller: "UserDeleteProfile handler"}

	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	user := models.User{ID: userID}
	if err := user.Delete(); err != nil {
		utils.RenderInternalServerErr(w)
		logger.Fatalln(err)
	}

	utils.DeleteTokenCookie(w)
	utils.SetPopupCookie(w, "Your account deleted successfully.")
	w.Header().Set("HX-Redirect", "/")
}
