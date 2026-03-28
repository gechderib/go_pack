package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var idCounter = 1

func CreateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println("the request id is: ", id)
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Username == "" {
		http.Error(w, "Missing required fields: username and email", http.StatusBadRequest)
		return
	}

	user.Id = idCounter
	idCounter++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Println("the request id is: ", id)

	// Here you would typically fetch the user from a database using the ID
	// For demonstration, we'll just return a dummy user
	user := User{
		Id:       1,
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
