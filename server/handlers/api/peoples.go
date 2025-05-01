package api

import (
	"effective/db"
	"effective/models"
	requestsmanager "effective/requests-manager"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Summary Добавить нового человека
// @Description Добавляет нового человека в базу данных с автоматическим определением возраста, пола и национальности
// @Tags people
// @Accept json
// @Produce json
// @Param person body models.PeopleFromRequest true "Данные человека"
// @Success 200 {object} models.People
// @Failure 400 {object} map[string]string
// @Failure 415 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /people [post]
func AddPeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	log.Printf("Получен запрос на добавление нового человека")
	var err error

	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("Ошибка: неверный Content-Type")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(map[string]string{"error": "Content-Type must be application/json"})
		return
	}

	var peopleRequest models.PeopleFromRequest
	if err := json.NewDecoder(r.Body).Decode(&peopleRequest); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Получение данных о человеке: %s %s", peopleRequest.Name, peopleRequest.Surname)
	people, err := requestsmanager.GetPeopleData(peopleRequest)
	if err != nil {
		log.Printf("Ошибка получения данных: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Добавление человека в базу данных")
	err = manager.AddPeople(people)
	if err != nil {
		log.Printf("Ошибка добавления в базу данных: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Человек успешно добавлен, ID: %d", people.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// @Summary Удалить человека
// @Description Удаляет человека из базы данных по ID
// @Tags people
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [delete]
func DeletePeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Удаление человека с ID: %d", id)
	err = manager.DeletePeople(id)
	if err != nil {
		log.Printf("Ошибка удаления: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Человек с ID %d успешно удален", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func UpdatePeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Обновление данных человека с ID: %d", id)
	var peopleRequest models.PeopleFromRequest
	if err := json.NewDecoder(r.Body).Decode(&peopleRequest); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	name, err := manager.GetPeopleName(id)
	if err != nil {
		log.Printf("Ошибка получения имени: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if peopleRequest.Name != name {
		log.Printf("Изменение имени, получение новых данных")
		people, err := requestsmanager.GetPeopleData(peopleRequest)
		if err != nil {
			log.Printf("Ошибка получения новых данных: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		err = manager.FullUpdatePeople(people, id)
		if err != nil {
			log.Printf("Ошибка полного обновления: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		log.Printf("Данные успешно обновлены")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(people)
		return
	}

	err = manager.UpdatePeople(&peopleRequest, id)
	if err != nil {
		log.Printf("Ошибка обновления: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Данные успешно обновлены")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peopleRequest)
}

func GetAllPeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	log.Printf("Получение списка людей")
	query := r.URL.Query()

	ageFrom, _ := strconv.Atoi(query.Get("age_from"))
	ageTo, _ := strconv.Atoi(query.Get("age_to"))
	gender := query.Get("gender")
	if gender != "" && gender != "male" && gender != "female" {
		log.Printf("Ошибка: неверное значение пола: %s", gender)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid gender value"})
		return
	}
	country := query.Get("country")

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	log.Printf("Параметры запроса: страница=%d, лимит=%d, возраст от=%d до=%d, пол=%s, страна=%s",
		page, limit, ageFrom, ageTo, gender, country)

	people, err := manager.GetAllPeople(ageFrom, ageTo, gender, country, offset, limit)
	if err != nil {
		log.Printf("Ошибка получения списка: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Успешно получено %d записей", len(people))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"page": page, "data": people})
}

func GetPeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	log.Printf("Получение данных о человеке")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Ошибка преобразования ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	people, err := manager.GetPeople(id)
	if err != nil {
		log.Printf("Ошибка получения данных: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Успешно получено данные о человеке")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}
