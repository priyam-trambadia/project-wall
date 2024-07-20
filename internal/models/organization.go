package models

import (
	"database/sql"

	"github.com/priyam-trambadia/project-wall/internal/logger"
)

type Organization struct {
	ID       int64  `json:"id"`
	Hostname string `json:"hostname"`
}

func (organization *Organization) Insert() error {
	logger := logger.Logger{Caller: "Organization::Insert model"}

	query := `
		INSERT INTO organizations (hostname)
		VALUES ($1)
		RETURNING id;
	`
	err := database.QueryRow(query, organization.Hostname).Scan(&organization.ID)
	err = logger.AppendError(err)
	return err
}

func (organization *Organization) Update() {}
func (organization *Organization) Get()    {}
func (organization *Organization) Delete() {}

// utilities

func GetOrCreateOrganizationID(hostname string) (int64, error) {
	logger := logger.Logger{Caller: "GetOrCreateOrganization model"}

	query := `
		SELECT id
		FROM organizations
		WHERE hostname = $1;
	`
	var id int64
	err := database.QueryRow(query, hostname).Scan(&id)
	if err == sql.ErrNoRows {
		organization := Organization{Hostname: hostname}
		if err := organization.Insert(); err != nil {
			err = logger.AppendError(err)
			return 0, err
		}
		id = organization.ID

	} else if err != nil {
		err = logger.AppendError(err)
		return 0, err
	}

	return id, nil
}
