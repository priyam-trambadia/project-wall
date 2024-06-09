package main

import "flag"

func initConstants(cfg *config) {

	flag.IntVar(&cfg.server.PORT, "port", 8000, "Server port")
	flag.StringVar(&cfg.database.URL, "database_url", "", "PostgreSQL URL")
	flag.Parse()
}
