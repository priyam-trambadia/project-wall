package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database_config struct {
	URL string
}

func (db_cfg Database_config) ConnectDB() *sql.DB {

	db, err := sql.Open("postgres", db_cfg.URL)
	if err != nil {
		log.Fatalln("[-] Fail to create database connection")
	}
	log.Println("[+] Database connected")

	return db
}
