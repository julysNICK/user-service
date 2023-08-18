package main

import (
	"fmt"
	"net/http"
	"strconv"
	"user/data"

	"github.com/go-chi/chi"
)

func readIDParam(r *http.Request) int {
	id := chi.URLParam(r, "id")

	convInt, _ := strconv.Atoi(id)

	return convInt
}

func (app *Config) getUser(w http.ResponseWriter, r *http.Request) {
	// Get a user by ID

	id := readIDParam(r)

	user, err := app.Models.User.GetById(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User retrieved",
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) getUsers(w http.ResponseWriter, r *http.Request) {
	// Get all users

	users, err := app.Models.User.GetAll()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Users retrieved",
		Data:    users,
	}

	app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) createUser(w http.ResponseWriter, r *http.Request) {
	// Create a new user

	var requestPayload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
		Active    int    `json:"active"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println(requestPayload)

	userCreate := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  requestPayload.Password,
	}


	user, err := app.Models.User.Insert(userCreate)


	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User created",
		Data:    user,
	}

	app.writeJSON(w, http.StatusCreated, payload)

}
