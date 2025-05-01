package models

// Country модель данных о стране
// @Description Модель данных о стране с вероятностью
type Country struct {
	Code        string  `json:"country_id" example:"RU"`
	Probability float64 `json:"probability" example:"0.95"`
}

// CountryResponse модель ответа сервиса определения национальности
// @Description Модель ответа от сервиса определения национальности
type CountryResponse struct {
	Count     int       `json:"count" example:"100"`
	Name      string    `json:"name" example:"Ivan"`
	Countries []Country `json:"country"`
}
