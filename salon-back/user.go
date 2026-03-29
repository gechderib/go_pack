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
	Password string `json:"password"` // Password should not be included in JSON responses
}

type CreateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password,omitempty"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (r CreateUserRequest) Validate() map[string]string {

	errors := make(map[string]string)

	if r.Username == nil {
		errors["username"] = "Username is required"
	} else if *r.Username == "" {
		errors["username"] = "Username can't be empty"
	}
	if r.Email == nil {
		errors["email"] = "Email is required"
	} else if *r.Email == "" {
		errors["email"] = "Email can't be empty"
	}
	return errors
}

func (r CreateUserRequest) ToModel(id int) User {
	return User{
		Id:       id,
		Username: *r.Username,
		Email:    *r.Email,
		Password: *r.Password,
	}
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}

func JSONMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		next(w, r)
	}
}

func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

var idCounter = 1

func CreateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println("the request id is: ", id)

	// var user User
	// err := json.NewDecoder(r.Body).Decode(&user)

	var req CreateUserRequest
	if err := DecodeJSON(r, &req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	errors := req.Validate()
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errors,
		})
		return
	}

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// if user.Email == "" || user.Username == "" {
	// 	http.Error(w, "Missing required fields: username and email", http.StatusBadRequest)
	// 	return
	// }
	user := req.ToModel(idCounter)
	user.Id = idCounter
	idCounter++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToUserResponse(user))

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
