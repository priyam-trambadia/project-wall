package main

import (
	"flag"

	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
)

func (cfg *config) initConfig() {

	flag.IntVar(&cfg.port, "port", 4000, "Server port")
	flag.StringVar(&cfg.database.url, "database_url", "", "PostgreSQL URL")
	flag.Parse()

	jwt.LoadConfig()
}
