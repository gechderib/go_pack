package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
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
		errors["username"] = "Username cannot be empty"
	}

	if r.Email == nil {
		errors["email"] = "Email is required"
	} else if *r.Email == "" {
		errors["email"] = "Email cannot be empty"
	}

	if r.Password == nil {
		errors["password"] = "Password is required"
	} else if *r.Password == "" {
		errors["password"] = "Password cannot be empty"
	}

	return errors
}

func (r CreateUserRequest) ToModel(id int) User {
	return User{
		ID:       id,
		Username: *r.Username,
		Email:    *r.Email,
		Password: *r.Password,
	}
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

var idCounter = 1

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	error := json.NewDecoder(r.Body).Decode(&req)

	if error != nil {
		http.Error(w, "Invalid request body: "+error.Error(), http.StatusBadRequest)
		return
	}

	errors := req.Validate()
	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	user := req.ToModel(idCounter)
	idCounter++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToUserResponse(user))

}
