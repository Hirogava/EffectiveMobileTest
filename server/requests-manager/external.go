package requestsmanager

import (
	"effective/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

func Request(client *http.Client, fullUrl string) ([]byte, error) {
	log.Printf("Отправка запроса к %s", fullUrl)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Printf("Ошибка создания запроса: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Неожиданный код ответа: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела ответа: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Успешный ответ от %s", fullUrl)
	return body, nil
}

func GetAge(personName string) (int, error) {
	log.Printf("Получение возраста для имени: %s", personName)
	var ageResponse models.PeopleAge

	baseUrl := models.AgifyURL
	params := url.Values{}
	params.Add("name", personName)
	fullUrl := baseUrl + "?" + params.Encode()

	body, err := Request(defaultClient, fullUrl)
	if err != nil {
		log.Printf("Ошибка получения возраста: %v", err)
		return 0, fmt.Errorf("failed to get age: %w", err)
	}

	if err := json.Unmarshal(body, &ageResponse); err != nil {
		log.Printf("Ошибка парсинга JSON ответа: %v", err)
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("Получен возраст %d для имени %s", ageResponse.Age, personName)
	return ageResponse.Age, nil
}

func GetGender(personName string) (string, error) {
	log.Printf("Получение пола для имени: %s", personName)
	var genderResponse models.PeopleGender
	male := models.GenderMale
	female := models.GenderFemale

	baseUrl := models.GenderizeURL
	params := url.Values{}
	params.Add("name", personName)
	fullUrl := baseUrl + "?" + params.Encode()

	body, err := Request(defaultClient, fullUrl)
	if err != nil {
		log.Printf("Ошибка получения пола: %v", err)
		return "", fmt.Errorf("failed to get gender: %w", err)
	}

	if err := json.Unmarshal(body, &genderResponse); err != nil {
		log.Printf("Ошибка парсинга JSON ответа: %v", err)
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	switch genderResponse.Gender {
	case male, female:
		log.Printf("Получен пол %s для имени %s", genderResponse.Gender, personName)
		return genderResponse.Gender, nil
	default:
		log.Printf("Неизвестный пол: %s", genderResponse.Gender)
		return "", fmt.Errorf("unknown gender: %s", genderResponse.Gender)
	}
}

func GetNationality(personName string) ([]models.Country, error) {
	log.Printf("Получение национальности для имени: %s", personName)
	var nationalityResponse models.CountryResponse

	baseUrl := models.NationalizeURL
	params := url.Values{}
	params.Add("name", personName)
	fullUrl := baseUrl + "?" + params.Encode()

	body, err := Request(defaultClient, fullUrl)
	if err != nil {
		log.Printf("Ошибка получения национальности: %v", err)
		return nil, fmt.Errorf("failed to get nationality: %w", err)
	}

	if err := json.Unmarshal(body, &nationalityResponse); err != nil {
		log.Printf("Ошибка парсинга JSON ответа: %v", err)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("Получено %d стран для имени %s", len(nationalityResponse.Countries), personName)
	return nationalityResponse.Countries, nil
}

func GetPeopleData(request models.PeopleFromRequest) (*models.People, error) {
	log.Printf("Получение данных о человеке: %s %s", request.Name, request.Surname)
	var err error
	var people models.People

	people.Age, err = GetAge(request.Name)
	if err != nil {
		log.Printf("Ошибка получения возраста: %v", err)
		return nil, fmt.Errorf("failed to get age: %w", err)
	}

	people.Gender, err = GetGender(request.Name)
	if err != nil {
		log.Printf("Ошибка получения пола: %v", err)
		return nil, fmt.Errorf("failed to get gender: %w", err)
	}

	people.Countries, err = GetNationality(request.Name)
	if err != nil {
		log.Printf("Ошибка получения национальности: %v", err)
		return nil, fmt.Errorf("failed to get nationality: %w", err)
	}

	people.Name = request.Name
	people.Surname = request.Surname
	people.Patronymic = request.Patronymic

	log.Printf("Успешно получены все данные о человеке: %s %s", request.Name, request.Surname)
	return &people, nil
}