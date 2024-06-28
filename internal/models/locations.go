package models

//
// Data source : https://simplemaps.com/data/world-cities
//

type Location struct {
	ID        int64  `json:"id"`
	City      string `json:"city"`
	CityASCII string `json:"city_ascii"`
	CityAlt   string `json:"city_alt"`
	AdminName string `json:"admin_name"`
	Country   string `json:"country"`
}

func (location *Location) Insert() error {
	query := `
		INSERT INTO locations (id, city, city_ascii, city_alt, admin_name, country)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	args := []interface{}{
		location.ID,
		location.City,
		location.CityASCII,
		location.CityAlt,
		location.AdminName,
		location.Country,
	}

	_, err := database.Exec(query, args)
	return err
}

func (location *Location) Update() {}
func (location *Location) Get()    {}
func (location *Location) Delete() {}

// utilities

func FindLocationsWithFullTextSearch(cityName string) ([]Location, error) {
	query := `
		SELECT id, city, city_ascii, city_alt, admin_name, country
		FROM locations
		WHERE (to_tsvector('simple', city) @@ plainto_tsquery('simple', $1) 
			OR to_tsvector('simple', city_ascii) @@ plainto_tsquery('simple', $1) 
			OR to_tsvector('simple', city_alt) @@ plainto_tsquery('simple', $1) 
			OR $1 = '');
	`
	if rows, err := database.Query(query, cityName); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		locations := make([]Location, 0)

		for rows.Next() {
			var location Location

			err := rows.Scan(
				&location.ID,
				&location.City,
				&location.CityASCII,
				&location.CityAlt,
				&location.AdminName,
				&location.Country,
			)

			if err != nil {
				return nil, err
			} else {
				locations = append(locations, location)
			}
		}

		return locations, nil
	}

}
