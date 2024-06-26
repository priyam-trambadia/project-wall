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

	srv.mux.HandleFunc("GET /user/register", handlers.UserRegister)
	srv.mux.HandleFunc("POST /user/register", handlers.UserRegisterPOST)

	srv.mux.HandleFunc("GET /user/login", handlers.UserLogin)
	srv.mux.HandleFunc("POST /user/login", handlers.UserLoginPOST)

	// srv.mux.HandleFunc("PUT /user/update", >>handlerHere)
	srv.mux.HandleFunc("GET /user/logout", handlers.UserLogout)

	// srv.mux.HandleFunc("GET /user/profile", >>handlerHere)

	srv.mux.HandleFunc(
		"GET /project/add",
		middlewares.AuthenticationRequired(handlers.ProjectAdd),
	)
	srv.mux.HandleFunc(
		"POST /project/add",
		middlewares.AuthenticationRequired(handlers.ProjectAddPOST),
	)

	// srv.mux.HandleFunc("PUT /project/update/{project_id}", >>handlerHere)
	// srv.mux.HandleFunc("DELETE /project/remove/{project_id}", >>handlerHere)
	// srv.mux.HandleFunc("PATCH /project/bookmark/{project_id}", >>handlerHere)

	// srv.mux.HandleFunc("GET /project/search", >>handlerHere)
	srv.mux.HandleFunc("GET /project/tag/search", handlers.ProjectTagSearch)
	srv.mux.HandleFunc("GET /project/language/search", handlers.ProjectLanguageSearch)
}
