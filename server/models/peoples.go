package models

// PeopleFromRequest модель для создания/обновления человека
// @Description Модель для создания или обновления данных о человеке
type PeopleFromRequest struct {
	Name       string  `json:"name" example:"Ivan"`
	Surname    string  `json:"surname" example:"Ivanov"`
	Patronymic *string `json:"patronymic,omitempty" example:"Ivanovich"`
}

// People полная модель данных о человеке
// @Description Полная модель данных о человеке, включая автоматически определенные данные
type People struct {
	ID         int       `json:"id" example:"1"`
	Name       string    `json:"name" example:"Ivan"`
	Surname    string    `json:"surname" example:"Ivanov"`
	Patronymic *string   `json:"patronymic,omitempty" example:"Ivanovich"`
	Age        int       `json:"age" example:"30"`
	Gender     string    `json:"gender" example:"male"`
	Countries  []Country `json:"countries"`
}

const (
	GenderMale   = "male"
	GenderFemale = "female"
)

// PeopleAge модель ответа сервиса определения возраста
// @Description Модель ответа от сервиса определения возраста
type PeopleAge struct {
	Count int    `json:"count" example:"100"`
	Name  string `json:"name" example:"Ivan"`
	Age   int    `json:"age" example:"30"`
}

// PeopleGender модель ответа сервиса определения пола
// @Description Модель ответа от сервиса определения пола
type PeopleGender struct {
	Count  int    `json:"count" example:"100"`
	Name   string `json:"name" example:"Ivan"`
	Gender string `json:"gender" example:"male"`
}