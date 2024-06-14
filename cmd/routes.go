package main

import "github.com/priyam-trambadia/project-wall/api/handlers"

func (srv *server) addRoutes() {

	srv.mux.HandleFunc("GET /{$}", handlers.Root)

	srv.mux.HandleFunc("GET /user/register", handlers.UserRegister)
	srv.mux.HandleFunc("POST /user/register", handlers.UserRegisterPOST)

	srv.mux.HandleFunc("GET /user/login", handlers.UserLogin)
	srv.mux.HandleFunc("POST /user/login", handlers.UserLoginPOST)

	srv.mux.HandleFunc("PUT /user/update", "[[[ add handler here ]]]")
	srv.mux.HandleFunc("POST /user/logout", handlers.UserLogoutPOST)

	srv.mux.HandleFunc("GET /project/add", "[[[ add handler here ]]]")
	srv.mux.HandleFunc("POST /project/add", "[[[ add handler here ]]]")

	srv.mux.HandleFunc("PUT /project/update/{project_id}", "[[[ add handler here ]]]")
	srv.mux.HandleFunc("DELETE /project/remove/{project_id}", "[[[ add handler here ]]]")
	srv.mux.HandleFunc("PATCH /project/bookmark/{project_id}", "[[[ add handler here ]]]")

	srv.mux.HandleFunc("GET /project/trending/tags", "[[[ add handler here ]]]")
	srv.mux.HandleFunc("GET /project/trending/languages", "[[[ add handler here ]]]")

}
