package main

import (
	"log"
	"net/http"
	"os"

	"github.com/darkgooddack/bookvault-api/config"
	"github.com/darkgooddack/bookvault-api/db"
	"github.com/darkgooddack/bookvault-api/handlers"
	"github.com/darkgooddack/bookvault-api/middleware"
	"github.com/darkgooddack/bookvault-api/models"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db.InitDB(cfg)

	if err := db.DB.AutoMigrate(&models.User{}, &models.Book{}); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	middleware.InitFromConfig(cfg)
	handlers.InitAuthHandler([]byte(cfg.JWTSecret))

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	books := r.PathPrefix("/books").Subrouter()
	books.Use(middleware.AuthMiddleware)

	books.HandleFunc("", handlers.GetBooks).Methods("GET")
	books.HandleFunc("", handlers.CreateBook).Methods("POST")
	books.HandleFunc("/{id}", handlers.GetBookByID).Methods("GET")
	books.HandleFunc("/{id}", handlers.UpdateBook).Methods("PUT")
	books.HandleFunc("/{id}", handlers.DeleteBook).Methods("DELETE")

	port := cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	addr := ":" + port

	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
