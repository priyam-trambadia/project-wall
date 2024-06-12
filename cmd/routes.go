package main

import "github.com/priyam-trambadia/project-wall/api/handlers"

func (srv *server) addRoutes() {

	srv.mux.HandleFunc("GET /{$}", handlers.Root)

	srv.mux.HandleFunc("GET /user/login", handlers.UserLogin)
	srv.mux.HandleFunc("POST /user/login", handlers.UserLoginPOST)

	srv.mux.HandleFunc("GET /user/register", handlers.UserRegister)
	srv.mux.HandleFunc("POST /user/register", handlers.UserRegisterPOST)
}
