package main

import (
	"net/http"

	"github.com/priyam-trambadia/project-wall/api/handlers"
	"github.com/priyam-trambadia/project-wall/api/middlewares"
)

func (srv *server) addRoutes() {

	srv.mux.HandleFunc("GET /{$}", handlers.Root)

	fileServer := http.FileServer(http.Dir("./web/static"))
	srv.mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// user routes
	srv.mux.HandleFunc("GET /user/register", handlers.UserRegister)
	srv.mux.HandleFunc("POST /user/register", handlers.UserRegisterPOST)
	srv.mux.HandleFunc("GET /user/login", handlers.UserLogin)
	srv.mux.HandleFunc("POST /user/login", handlers.UserLoginPOST)
	srv.mux.HandleFunc("GET /user/activate", handlers.UserActivate)
	srv.mux.HandleFunc("GET /user/forgot-password", handlers.UserForgotPassword)
	srv.mux.HandleFunc("POST /user/forgot-password", handlers.UserForgotPasswordPOST)
	srv.mux.HandleFunc("GET /user/password/reset", handlers.UserPasswordReset)
	srv.mux.HandleFunc("POST /user/password/reset", handlers.UserPasswordResetPOST)
	srv.mux.HandleFunc(
		"GET /user/logout",
		middlewares.UserAuthenticationRequired(handlers.UserLogout),
	)
	srv.mux.HandleFunc("GET /user/{user_id}/avatar", handlers.UserAvatar)
	srv.mux.HandleFunc(
		"GET /user/{user_id}",
		middlewares.ValidatePathValueUserID(handlers.UserGetProfile),
	)
	srv.mux.HandleFunc(
		"PUT /user/{user_id}",
		middlewares.UserAuthenticationRequired(
			middlewares.ValidatePathValueUserID(
				middlewares.AuthorizeUserProfileAction(handlers.UserUpdateProfile),
			),
		),
	)
	srv.mux.HandleFunc(
		"DELETE /user/{user_id}",
		middlewares.UserAuthenticationRequired(
			middlewares.ValidatePathValueUserID(
				middlewares.AuthorizeUserProfileAction(handlers.UserDeleteProfile),
			),
		),
	)

	// project routes
	srv.mux.HandleFunc(
		"GET /project/create",
		middlewares.UserAuthenticationRequired(handlers.ProjectCreate),
	)
	srv.mux.HandleFunc(
		"POST /project/create",
		middlewares.UserAuthenticationRequired(handlers.ProjectCreatePOST),
	)
	srv.mux.HandleFunc(
		"PUT /project/{project_id}",
		middlewares.UserAuthenticationRequired(
			middlewares.ValidatePathValueProjectID(
				middlewares.AuthorizeProjectAction(handlers.ProjectUpdate),
			),
		),
	)
	srv.mux.HandleFunc(
		"DELETE /project/{project_id}",
		middlewares.UserAuthenticationRequired(
			middlewares.ValidatePathValueProjectID(
				middlewares.AuthorizeProjectAction(handlers.ProjectDelete),
			),
		),
	)
	srv.mux.HandleFunc(
		"PATCH /project/{project_id}/toggle-bookmark",
		middlewares.UserAuthenticationRequired(
			middlewares.ValidatePathValueProjectID(handlers.ProjectToggleBookmark),
		),
	)

	srv.mux.HandleFunc("GET /project/search", handlers.ProjectSearch)
	srv.mux.HandleFunc("GET /project/tag/search", handlers.ProjectTagSearch)
	srv.mux.HandleFunc("GET /project/language/search", handlers.ProjectLanguageSearch)
}
