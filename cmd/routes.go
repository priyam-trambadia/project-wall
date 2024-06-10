package main

import "github.com/priyam-trambadia/project-wall/api/handlers"

func (app *application) addRoutes() {

	app.mux.HandleFunc("GET /{$}", handlers.Root)

	app.mux.HandleFunc("GET /user/login", handlers.Login)
	app.mux.HandleFunc("POST /user/login", handlers.LoginPOST)

	app.mux.HandleFunc("GET /user/register", handlers.Register)
	app.mux.HandleFunc("POST /user/register", handlers.RegisterPOST)
}
