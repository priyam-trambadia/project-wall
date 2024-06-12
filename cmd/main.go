package main

import (
	"github.com/priyam-trambadia/project-wall/internal/models"
)

type config struct {
	port     int
	database struct {
		url string
	}
}

func main() {
	var cfg config
	cfg.initConfig()

	db := cfg.connectDatabase()
	defer db.Close()
	models.SetDatabaseVar(db)

	cfg.startServer()
}
