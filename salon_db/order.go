package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Order struct {
	BaseModel
	UserID      uint   `json:"user_id" gorm:"not null"`
	Description string `json:"description" gorm:"not null; size:255"`
}

type CreateOrderRequest struct {
	UserID      *uint   `json:"user_id"`
	Description *string `json:"description"`
}

type OrderResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToOrderResponse(order Order) OrderResponse {
	return OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Description: order.Description,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

func ToOrderResponses(orders []Order) []OrderResponse {
	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = ToOrderResponse(order)
	}
	return responses
}

func (r CreateOrderRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if r.UserID == nil {
		errors["user_id"] = "User ID is required"
	}

	if r.Description == nil {
		errors["description"] = "Description is required"
	} else if *r.Description == "" {
		errors["description"] = "Description cannot be empty"
	}

	return errors
}

func (r CreateOrderRequest) ToModel() Order {
	return Order{
		UserID:      *r.UserID,
		Description: *r.Description,
	}
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	errors := req.Validate()
	if len(errors) > 0 {
		JSONResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Errors:  errors,
		})
		return
	}

	order := req.ToModel()
	if err := db.Create(&order).Error; err != nil {
		JSONResponse(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create order",
			Errors:  err.Error(),
		})
		return
	}

	JSONResponse(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "Order created successfully",
		Data:    ToOrderResponse(order),
	})
}
