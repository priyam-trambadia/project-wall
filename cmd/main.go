package main

import (
	"github.com/priyam-trambadia/project-wall/api/utils/jwt"
	"github.com/priyam-trambadia/project-wall/internal/mailer"
	"github.com/priyam-trambadia/project-wall/internal/models"
)

func main() {

	jwt.LoadConfig()
	mailer.SetupMailer()

	var cfg config
	cfg.initConfig()

	db := cfg.connectDatabase()
	defer db.Close()
	models.SetDatabaseVar(db)

	cfg.startServer()
}
