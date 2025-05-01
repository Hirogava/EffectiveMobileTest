package db

import (
	"database/sql"
	"effective/models"
)

func (manager *Manager) AddPeopleCountries(tx *sql.Tx, people *models.People, peopleId int) error {
	for _, country := range people.Countries {
		query := `INSERT INTO person_countries(person_id, country_id, probability) VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, peopleId, country.Code, country.Probability)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (manager *Manager) DeletePeopleCountries(tx *sql.Tx, personId int) error {
	_, err := tx.Exec("DELETE FROM person_countries WHERE person_id = $1", personId)
	return err
}