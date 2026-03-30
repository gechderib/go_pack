package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	db_postgres_url := "host=localhost user=myuser password=mypass dbname=mydb port=5436"
	db, err = gorm.Open(postgres.Open(db_postgres_url), &gorm.Config{})

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
	userRouter.Delete("/users/{id}", DeleteUser)

	http.ListenAndServe(":8080", userRouter)
}
