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
	apiRout := r.PathPrefix("/api").Subrouter()

	apiRout.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		api.AddPeople(w, r, manager)
	}).Methods(http.MethodPost)

	apiRout.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.GetPeople(w, r, manager)
	}).Methods(http.MethodGet)

	apiRout.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.DeletePeople(w, r, manager)
	}).Methods(http.MethodDelete)

	apiRout.HandleFunc("/people/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.UpdatePeople(w, r, manager)
	}).Methods(http.MethodPut)

	apiRout.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		api.GetAllPeople(w, r, manager)
	}).Methods(http.MethodGet)
}