package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type User struct {
	BaseModel
	FullName string  `json:"full_name" gorm:"not null; size:255"`
	Username string  `json:"username" gorm:"not null; size:255; unique; uniqueIndex"`
	Email    string  `json:"email" gorm:"not null; size:255; unique; uniqueIndex"`
	Phone    string  `json:"phone" gorm:"size:20"`
	Password string  `json:"password" gorm:"not null; size:255"`
	Orders   []Order `json:"orders" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type CreateUserRequest struct {
	FullName *string `json:"full_name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Phone    *string `json:"phone"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	if r.Phone == nil {
		errors["phone"] = "Phone is required"
	} else if *r.Phone == "" {
		errors["phone"] = "Phone cannot be empty"
	}

	return errors
}

func (r CreateUserRequest) ToModel() User {
	return User{
		FullName: *r.FullName,
		Username: *r.Username,
		Email:    *r.Email,
		Phone:    *r.Phone,
		Password: *r.Password,
	}
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}

type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
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

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "" {
		pageStr = "1"
	}
	if limitStr == "" {
		limitStr = "10"
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	time.Sleep(4 * time.Second) // Simulate a long-running operation
	var users []User
	// result := db.Find(&users)
	result := db.Offset((page - 1) * limit).Limit(limit).Find(&users)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve users: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(ToUserResponses(users))
	JSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    ToUserResponses(users),
		Pagination: &Pagination{
			Total: int(result.RowsAffected),
			Page:  page,
			Limit: limit,
		},
	})
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result := db.Delete(&User{}, id)
	if result.Error != nil {
		JSONResponse(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete user",
			Errors:  result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		JSONResponse(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}
	JSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "User deleted	successfully",
	})
}
