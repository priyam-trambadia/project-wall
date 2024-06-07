package main

import (
	"github.com/priyam-trambadia/project-wall/src/db"
)

func main() {
	db := db.ConnectDB()
	defer db.Close()
}
