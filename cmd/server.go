package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/priyam-trambadia/project-wall/api/middlewares"
)

type server struct {
	mux *http.ServeMux
}

func (cfg *config) startServer() {

	var srv server
	srv.mux = http.NewServeMux()
	srv.addRoutes()

	log.Println("[+] Start server listening on", cfg.port)

	addr := fmt.Sprintf(":%d", cfg.port)
	err := http.ListenAndServe(addr, middlewares.Authenticate(srv.mux))
	if err != nil {
		log.Fatalln("[-] Fail to create server")
	}
}
