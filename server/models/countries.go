package models

type Country struct {
    Code        string  `json:"country_id"`
    Probability float64 `json:"probability"`
}

type CountryResponse struct {
    Count     int       `json:"count"`
    Name      string    `json:"name"`
    Countries []Country `json:"country"`
}