package main

import (
	"net/http"
	"user/data"
)

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

	userCreate := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  requestPayload.Password,
		Active:    requestPayload.Active,
	}

	user, err := app.Models.User.Insert(userCreate)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, user)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	return

}
