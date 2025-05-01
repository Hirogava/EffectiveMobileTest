package db

import (
	"database/sql"
	"effective/models"
	"log"
)

func (manager *Manager) AddPeopleCountries(tx *sql.Tx, people *models.People, peopleId int) error {
	log.Printf("Добавление %d стран для человека ID: %d", len(people.Countries), peopleId)
	for i, country := range people.Countries {
		query := `INSERT INTO person_countries(person_id, country_id, probability) VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, peopleId, country.Code, country.Probability)
		if err != nil {
			log.Printf("Ошибка добавления страны %d: %v", i+1, err)
			return err
		}
	}

	log.Printf("Успешно добавлены все страны для человека ID: %d", peopleId)
	return nil
}

func (manager *Manager) DeletePeopleCountries(tx *sql.Tx, personId int) error {
	log.Printf("Удаление стран для человека ID: %d", personId)
	_, err := tx.Exec("DELETE FROM person_countries WHERE person_id = $1", personId)
	if err != nil {
		log.Printf("Ошибка удаления стран: %v", err)
	} else {
		log.Printf("Страны успешно удалены для человека ID: %d", personId)
	}
	return err
}

func (manager *Manager) GetPeopleCountries(id int) ([]models.Country, error) {
	log.Printf("Получение стран для человека ID: %d", id)
	rows, err := manager.Conn.Query("SELECT country_id, probability FROM person_countries WHERE person_id = $1", id)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return nil, err
	}
	defer rows.Close()

	countries := []models.Country{}
	for rows.Next() {
		var country models.Country
		err = rows.Scan(&country.Code, &country.Probability)
		if err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, err
		}
		countries = append(countries, country)
	}

	log.Printf("Успешно получено %d стран для человека ID: %d", len(countries), id)
	return countries, nil
}