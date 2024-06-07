package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("[-] Error in loding .env for database")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("[-] Fail to create database connection")
	}

	log.Println("[+] Database connected")

	return db
}
