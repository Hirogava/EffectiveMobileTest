package requestsmanager

import (
	"effective/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var defaultClient = &http.Client{
    Timeout: 10 * time.Second,
}

func Request(client *http.Client, fullUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func GetAge(personName string) (int, error) {
	var ageResponse models.PeopleAge

	baseUrl := models.AgifyURL
	params := url.Values{}
	params.Add("name", personName)
	fullUrl := baseUrl + "?" + params.Encode()

	body, err := Request(defaultClient, fullUrl)
	if err != nil {
		return 0, fmt.Errorf("failed to get age: %w", err)
	}

	if err := json.Unmarshal(body, &ageResponse); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return ageResponse.Age, nil
}

func GetGender(personName string) (string, error) {
	var genderResponse models.PeopleGender
	male := models.GenderMale
	female := models.GenderFemale

	baseUrl := models.GenderizeURL
	params := url.Values{}
    params.Add("name", personName)
	fullUrl := baseUrl + "?" + params.Encode()

	body, err := Request(defaultClient, fullUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get gender: %w", err)
	}

	if err := json.Unmarshal(body, &genderResponse); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	switch genderResponse.Gender {
	case male, female:
		return genderResponse.Gender, nil
	default:
		return "", fmt.Errorf("unknown gender: %s", genderResponse.Gender)
	}
}

func GetNationality(personName string) ([]models.Country, error) {
	var nationalityResponse models.CountryResponse

	baseUrl := models.NationalizeURL
	params := url.Values{}
	params.Add("name", personName)
	fullUrl := baseUrl + "?" + params.Encode()

	body, err := Request(defaultClient, fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get nationality: %w", err)
	}

	if err := json.Unmarshal(body, &nationalityResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nationalityResponse.Countries, nil
}

func GetPeopleData(request models.PeopleFromRequest) (*models.People, error) {
	var err error
	var people models.People

	people.Age, err = GetAge(request.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get age: %w", err)
	}

	people.Gender, err = GetGender(request.Name)
	if err != nil {
        return nil, fmt.Errorf("failed to get gender: %w", err)
	}

	people.Countries, err = GetNationality(request.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get nationality: %w", err)
	}

	people.Name = request.Name
	people.Surname = request.Surname
	people.Patronymic = request.Patronymic

	return &people, nil
}