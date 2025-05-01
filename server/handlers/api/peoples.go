package api

import (
	"effective/db"
	"effective/models"
	"effective/requests-manager"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddPeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	var err error

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(map[string]string{"error": "Content-Type must be application/json"})
		return
	}

	var peopleRequest models.PeopleFromRequest

	if err := json.NewDecoder(r.Body).Decode(&peopleRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	people, err := requestsmanager.GetPeopleData(peopleRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = manager.AddPeople(people)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func DeletePeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = manager.DeletePeople(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func UpdatePeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
  	}

	var peopleRequest models.PeopleFromRequest
	if err := json.NewDecoder(r.Body).Decode(&peopleRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	name, err := manager.GetPeopleName(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if peopleRequest.Name != name {
		people, err := requestsmanager.GetPeopleData(peopleRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		err = manager.FullUpdatePeople(people, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(people)
		return
	}

	err = manager.UpdatePeople(&peopleRequest, id)
	if err != nil {
    	w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peopleRequest)
}

func GetPeople(w http.ResponseWriter, r *http.Request, manager *db.Manager) {
	query := r.URL.Query()
    
    ageFrom, _ := strconv.Atoi(query.Get("age_from"))
    ageTo, _ := strconv.Atoi(query.Get("age_to"))
    gender := query.Get("gender")
	if gender != "" && gender != "male" && gender != "female" {
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

	people, err := manager.GetAllPeople(ageFrom, ageTo, gender, country, offset, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"page" : page, "data": people}) 
}