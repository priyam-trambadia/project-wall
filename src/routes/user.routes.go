package routes

import "github.com/priyam-trambadia/project-wall/src/controllers"

func (r Route) AddUserRoute(base_route string) {

	r.Mux.HandleFunc("GET "+base_route+"/login", controllers.Login)
	r.Mux.HandleFunc("POST "+base_route+"/login", controllers.LoginPOST)

	r.Mux.HandleFunc("GET "+base_route+"/register", controllers.Register)
	r.Mux.HandleFunc("POST "+base_route+"/register", controllers.RegisterPOST)
}
