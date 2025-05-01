package models

type PeopleFromRequest struct {
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
}

type People struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic *string   `json:"patronymic,omitempty"`
	Age        int       `json:"age"`
	Gender     string    `json:"gender"`
	Countries  []Country `json:"countries"`
}

const (
	GenderMale   = "male"
	GenderFemale = "female"
)

type PeopleAge struct {
	Count      int       `json:"count"`
    Name       string    `json:"name"`
    Age        int       `json:"age"`
}

type PeopleGender struct {
    Count      int       `json:"count"`
    Name       string    `json:"name"`
    Gender     string    `json:"gender"`
}