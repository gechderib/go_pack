package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
}

func main() {
	initDB()

	userRouter := chi.NewRouter()
	userRouter.Post("/users", CreateUser)
	userRouter.Get("/users", GetUsers)
	userRouter.Get("/users/{id}", GetUserById)

	http.ListenAndServe(":8080", userRouter)
}
