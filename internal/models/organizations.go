package models

type Organization struct {
	ID       int64
	Hostname string
}

func (organization *Organization) Insert() error {
	query := `
		INSERT INTO organizations (hostname)
		VALUES ($1)
		RETURNING id;
	`
	return database.QueryRow(query, organization.Hostname).Scan(&organization.ID)
}

func (organization *Organization) Update() {}
func (organization *Organization) Get()    {}
func (organization *Organization) Delete() {}

// utilities

func GetOrganizationID(hostname string) (int64, error) {
	query := `
		SELECT id
		FROM organizations
		WHERE hostname = $1;
	`
	var id int64
	if err := database.QueryRow(query, hostname).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
