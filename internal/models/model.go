package models

import (
	"database/sql"
	"log"
)

var database *sql.DB

func SetDatabaseVar(db *sql.DB) {

	database = db

	err := database.Ping()
	if err != nil {
		log.Fatalln("[-] Error in setting model database variable")
	}

	log.Println("[+] Model database variable set successful")
}
