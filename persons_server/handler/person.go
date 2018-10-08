package handler

import (
	"encoding/json"
	"net/http"
	"testTask/persons_server/model"
	"testTask/persons_server/repository"
)

const HandlePath = "/persons/"

type personHandler struct {
	storage repository.PersonRepository
}

func NewPersonHandler(repository repository.PersonRepository) *personHandler {
	return &personHandler{repository}
}

func (ph *personHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == http.MethodPost { //if there is more then 2 clauses it's better to use switch-case operator
		ph.addPersonRecord(w, r)
	}
}

func (ph *personHandler) addPersonRecord(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var person model.Person
	err := decoder.Decode(&person)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	ph.storage.SavePerson(person)
	w.WriteHeader(http.StatusCreated)
}
