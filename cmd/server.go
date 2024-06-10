package main

import (
	"fmt"
	"log"
	"net/http"
)

type application struct {
	mux *http.ServeMux
}

func (cfg *config) startServer() {

	var app application
	app.mux = http.NewServeMux()
	app.addRoutes()

	log.Println("[+] Start server listening on", cfg.port)

	addr := fmt.Sprintf(":%d", cfg.port)
	err := http.ListenAndServe(addr, app.mux)
	if err != nil {
		log.Fatalln("[-] Fail to create server")
	}
}
