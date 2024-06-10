package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func (cfg *config) connectDatabase() *sql.DB {

	db, err := sql.Open("postgres", cfg.database.url)
	if err != nil {
		log.Fatalln("[-] Fail to create database connection")
	}
	log.Println("[+] Database connected")

	return db
}
