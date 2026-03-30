package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	FullName  string    `json:"full_name" gorm:"not null; size:255"`
	Username  string    `json:"username" gorm:"not null; size:255; unique; uniqueIndex"`
	Email     string    `json:"email" gorm:"not null; size:255; unique; uniqueIndex"`
	Password  string    `json:"password" gorm:"not null; size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type CreateUserRequest struct {
	FullName *string `json:"full_name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
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

	if r.FullName == nil {
		errors["full_name"] = "Full name is required"
	} else if *r.FullName == "" {
		errors["full_name"] = "Full name cannot be empty"
	}

	return errors
}

func (r CreateUserRequest) ToModel() User {
	return User{
		FullName: *r.FullName,
		Username: *r.Username,
		Email:    *r.Email,
		Password: *r.Password,
	}
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
	}
}

func ToUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSONResponse(w http.ResponseWriter, status int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	errors := req.Validate()
	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	user := req.ToModel()

	result := db.Create(&user)
	if result.Error != nil {
		JSONResponse(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  result.Error.Error(),
		})
		// http.Error(w, "Failed to create user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(ToUserResponse(user))

	JSONResponse(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    ToUserResponse(user),
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve users: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToUserResponses(users))

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToUserResponse(user))
}
