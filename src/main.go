package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/priyam-trambadia/project-wall/src/db"
)

type config struct {
	port     int
	database db.Database_config
}

func main() {
	var cfg config
	initConstants(&cfg)

	db := cfg.database.ConnectDB()
	defer db.Close()

	mux := http.NewServeMux()

	log.Println("[+] Start server listening on", cfg.port)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.port),
		mux)
	if err != nil {
		log.Fatalln("[-] Fail to create server")
	}

}
