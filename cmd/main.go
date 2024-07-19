package main

import (
	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
	"github.com/priyam-trambadia/project-wall/internal/mail"
	"github.com/priyam-trambadia/project-wall/internal/models"
)

func main() {

	jwt.LoadConfig()
	mail.SetupMailer()

	var cfg config
	cfg.initConfig()

	db := cfg.connectDatabase()
	defer db.Close()
	models.SetDatabaseVar(db)

	cfg.startServer()
}
