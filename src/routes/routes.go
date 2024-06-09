package routes

import "net/http"

type Route struct {
	Mux *http.ServeMux
}
