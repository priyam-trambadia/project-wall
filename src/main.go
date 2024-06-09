package main

import (
	"github.com/priyam-trambadia/project-wall/src/db"
	"github.com/priyam-trambadia/project-wall/src/server"
)

type config struct {
	database db.Database_config
	server   server.Server_config
}

func main() {
	var cfg config
	initConstants(&cfg)

	db := cfg.database.ConnectDB()
	defer db.Close()

	cfg.server.StartServer()
}
