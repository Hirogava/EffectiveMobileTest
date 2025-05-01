package routes

import (
	"effective/db"
	"effective/handlers/api"
	"net/http"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router, manager *db.Manager) {
	ApiRoutes(r, manager)
}

func ApiRoutes(r *mux.Router, manager *db.Manager) {
	r.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		api.AddPeople(w, r, manager)
	}).Methods(http.MethodPost)

	r.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.GetPeople(w, r, manager)
	}).Methods(http.MethodGet)

	r.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.DeletePeople(w, r, manager)
	}).Methods(http.MethodDelete)

	r.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.UpdatePeople(w, r, manager)
	}).Methods(http.MethodPut)

	r.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		api.GetAllPeople(w, r, manager)
	}).Methods(http.MethodGet)
}