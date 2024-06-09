package routes

import "github.com/priyam-trambadia/project-wall/src/controllers"

func (r Route) AddRootRoute(base_route string) {
	r.Mux.HandleFunc(base_route+"/{$}", controllers.Home)
}
