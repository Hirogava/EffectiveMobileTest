package db

import (
	"effective/models"
	"encoding/json"
	"fmt"
)

func (manager *Manager) AddPeople(people *models.People) error {
	tx, err := manager.Conn.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

	err = tx.QueryRow("INSERT INTO peoples (name, surname, patronymic, age, gender) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		people.Name, people.Surname, people.Patronymic, people.Age, people.Gender).Scan(&people.ID)
	if err != nil {
		return err
	}

	err = manager.AddPeopleCountries(tx, people, people.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (manager *Manager) DeletePeople(id int) error {
	_, err := manager.Conn.Exec("DELETE FROM peoples WHERE id=$1", id)
    return err
}

func (manager *Manager) UpdatePeople(people *models.PeopleFromRequest, peopleId int) error {
	_, err := manager.Conn.Exec("UPDATE peoples SET name=$1, surname=$2, patronymic=$3 WHERE id=$4",
		people.Name, people.Surname, people.Patronymic, peopleId)
	return err
}

func (manager *Manager) FullUpdatePeople(people *models.People, peopleId int) error {
	tx, err := manager.Conn.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

	err = manager.DeletePeopleCountries(tx, peopleId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE peoples SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5 WHERE id=$6",
	people.Name, people.Surname, people.Patronymic, people.Age, people.Gender, peopleId)

	err = manager.AddPeopleCountries(tx, people, peopleId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (manager *Manager) GetPeopleName(id int) (string, error) {
	var name string
	err := manager.Conn.QueryRow("SELECT name FROM peoples WHERE id=$1", id).Scan(&name)
	return name, err
}

func (manager *Manager) GetAllPeople(ageFrom, ageTo int, gender, country string, offset, limit int) ([]models.People, error) {
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
            return nil, fmt.Errorf("scan failed: %w", err)
        }
        
        if err := json.Unmarshal(countries, &people.Countries); err != nil {
            return nil, fmt.Errorf("failed to unmarshal countries: %w", err)
        }
        
        peoples = append(peoples, people)
    }
    
    return peoples, nil
}