package db

import (
	"effective/models"
	"encoding/json"
	"fmt"
	"log"
)

func (manager *Manager) AddPeople(people *models.People) error {
	log.Printf("Добавление нового человека: %s %s", people.Name, people.Surname)
	tx, err := manager.Conn.Begin()
	if err != nil {
		log.Printf("Ошибка начала транзакции: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			log.Printf("Откат транзакции из-за ошибки: %v", err)
			tx.Rollback()
		}
	}()

	err = tx.QueryRow("INSERT INTO peoples (name, surname, patronymic, age, gender) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		people.Name, people.Surname, people.Patronymic, people.Age, people.Gender).Scan(&people.ID)
	if err != nil {
		log.Printf("Ошибка добавления человека: %v", err)
		return err
	}

	log.Printf("Добавление стран для человека ID: %d", people.ID)
	err = manager.AddPeopleCountries(tx, people, people.ID)
	if err != nil {
		log.Printf("Ошибка добавления стран: %v", err)
		return err
	}

	log.Printf("Успешное добавление человека ID: %d", people.ID)
	return tx.Commit()
}

func (manager *Manager) DeletePeople(id int) error {
	log.Printf("Удаление человека с ID: %d", id)
	_, err := manager.Conn.Exec("DELETE FROM peoples WHERE id=$1", id)
	if err != nil {
		log.Printf("Ошибка удаления человека: %v", err)
	} else {
		log.Printf("Человек с ID %d успешно удален", id)
	}
	return err
}

func (manager *Manager) UpdatePeople(people *models.PeopleFromRequest, peopleId int) error {
	log.Printf("Обновление данных человека ID: %d", peopleId)
	_, err := manager.Conn.Exec("UPDATE peoples SET name=$1, surname=$2, patronymic=$3 WHERE id=$4",
		people.Name, people.Surname, people.Patronymic, peopleId)
	if err != nil {
		log.Printf("Ошибка обновления данных: %v", err)
	} else {
		log.Printf("Данные человека ID %d успешно обновлены", peopleId)
	}
	return err
}

func (manager *Manager) FullUpdatePeople(people *models.People, peopleId int) error {
	log.Printf("Полное обновление данных человека ID: %d", peopleId)
	tx, err := manager.Conn.Begin()
	if err != nil {
		log.Printf("Ошибка начала транзакции: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			log.Printf("Откат транзакции из-за ошибки: %v", err)
			tx.Rollback()
		}
	}()

	log.Printf("Удаление старых стран для человека ID: %d", peopleId)
	err = manager.DeletePeopleCountries(tx, peopleId)
	if err != nil {
		log.Printf("Ошибка удаления стран: %v", err)
		return err
	}

	_, err = tx.Exec("UPDATE peoples SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5 WHERE id=$6",
		people.Name, people.Surname, people.Patronymic, people.Age, people.Gender, peopleId)
	if err != nil {
		log.Printf("Ошибка обновления данных: %v", err)
		return err
	}

	log.Printf("Добавление новых стран для человека ID: %d", peopleId)
	err = manager.AddPeopleCountries(tx, people, peopleId)
	if err != nil {
		log.Printf("Ошибка добавления стран: %v", err)
		return err
	}

	log.Printf("Успешное полное обновление данных человека ID: %d", peopleId)
	return tx.Commit()
}

func (manager *Manager) GetPeopleName(id int) (string, error) {
	log.Printf("Получение имени человека с ID: %d", id)
	var name string
	err := manager.Conn.QueryRow("SELECT name FROM peoples WHERE id=$1", id).Scan(&name)
	if err != nil {
		log.Printf("Ошибка получения имени: %v", err)
	} else {
		log.Printf("Получено имя: %s", name)
	}
	return name, err
}

func (manager *Manager) GetAllPeople(ageFrom, ageTo int, gender, country string, offset, limit int) ([]models.People, error) {
	log.Printf("Получение списка людей с параметрами: возраст от=%d до=%d, пол=%s, страна=%s, страница=%d, лимит=%d",
		ageFrom, ageTo, gender, country, offset/limit+1, limit)
	query := `SELECT 
			p.id,
			p.name,
			p.surname,
			p.patronymic,
			p.age,
			p.gender,
			(
				SELECT json_agg(
					json_build_object(
						'country_id', pc.country_id,
						'probability', pc.probability
					)
				)
				FROM person_countries pc
				WHERE pc.person_id = p.id
			) AS countries
		FROM 
			peoples p
		WHERE 
			($1 = 0 OR p.age >= $1) AND
			($2 = 0 OR p.age <= $2) AND
			($3 = '' OR p.gender = $3::gender) AND
			(
				$4 = '' OR 
				EXISTS (
					SELECT 1 
					FROM person_countries pc 
					WHERE pc.person_id = p.id AND pc.country_id = $4
				)
			)
		ORDER BY 
			p.id
		LIMIT $5 OFFSET $6;`

	var peoples []models.People
	rows, err := manager.Conn.Query(query, ageFrom, ageTo, gender, country, limit, offset)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var people models.People
		var countries []byte

		err := rows.Scan(
			&people.ID,
			&people.Name,
			&people.Surname,
			&people.Patronymic,
			&people.Age,
			&people.Gender,
			&countries,
		)

		if err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if err := json.Unmarshal(countries, &people.Countries); err != nil {
			log.Printf("Ошибка разбора JSON стран: %v", err)
			return nil, fmt.Errorf("failed to unmarshal countries: %w", err)
		}

		peoples = append(peoples, people)
	}

	log.Printf("Успешно получено %d записей", len(peoples))
	return peoples, nil
}

func (manager *Manager) GetPeople(id int) (models.People, error) {
	log.Printf("Получение данных о человеке с ID: %d", id)
	var people models.People
	err := manager.Conn.QueryRow("SELECT * FROM peoples WHERE id=$1", id).Scan(&people)
	if err != nil {
		log.Printf("Ошибка получения данных: %v", err)
		return models.People{}, err
	}

	log.Printf("Получение стран для человека ID: %d", id)
	countries, err := manager.GetPeopleCountries(id)
	if err != nil {
		log.Printf("Ошибка получения стран: %v", err)
		return models.People{}, err
	}

	people.Countries = countries
	log.Printf("Успешно получены все данные о человеке ID: %d", id)
	return people, nil
}