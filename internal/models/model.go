package models

import (
	"database/sql"
	"fmt"
	"log"
)

var (
	database          *sql.DB
	ErrRecordNotFound = sql.ErrNoRows
)

func SetDatabaseVar(db *sql.DB) {

	database = db

	err := database.Ping()
	if err != nil {
		log.Fatalln("[-] Error in setting model database variable")
	}

	log.Println("[+] Model database variable set successful")
}

func ArraytoStringRoundBrackets(arr []int64) string {
	str := "("

	for _, ele := range arr {
		str += fmt.Sprintf("%d,", ele)
	}

	str += "0)"

	return str
}
