package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/priyam-trambadia/project-wall/api/middlewares"
)

type config struct {
	port     int
	database struct {
		url string
	}
}

type server struct {
	mux *http.ServeMux
}

func (cfg *config) initConfig() {

	flag.IntVar(&cfg.port, "port", 4000, "Server port")
	flag.StringVar(&cfg.database.url, "database_url", "", "PostgreSQL URL")
	flag.Parse()
}

func (cfg *config) connectDatabase() *sql.DB {

	db, err := sql.Open("postgres", cfg.database.url)
	if err != nil {
		log.Fatalln("[-] Fail to create database connection")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("[-] Fail to establish database connection")
	}
	log.Println("[+] Database connected successfully")

	return db
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
