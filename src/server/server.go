package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/priyam-trambadia/project-wall/src/routes"
)

type Server_config struct {
	PORT int
}

type Application struct {
	route routes.Route
}

func (server_cfg Server_config) StartServer() {

	var app Application

	app.route.Mux = http.NewServeMux()
	app.route.AddRootRoute("")
	app.route.AddUserRoute("/user")

	log.Println("[+] Start server listening on", server_cfg.PORT)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", server_cfg.PORT),
		app.route.Mux)
	if err != nil {
		log.Fatalln("[-] Fail to create server")
	}
}
