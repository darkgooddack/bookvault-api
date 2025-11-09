package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/darkgooddack/bookvault-api/db"
	"github.com/darkgooddack/bookvault-api/middleware"
	"github.com/darkgooddack/bookvault-api/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GET /books  (только свои книги)
func GetBooks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var books []models.Book
	if err := db.DB.Where("owner = ?", userID).Find(&books).Error; err != nil {
		http.Error(w, "failed to query books", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}

// POST /books
func CreateBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	b.ID = uuid.New().String()
	b.Owner = userID

	if err := db.DB.Create(&b).Error; err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

// GET /books/{id}
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	var b models.Book
	if err := db.DB.First(&b, "id = ?", id).Error; err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	if b.Owner != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(b)
}

// PUT /books/{id}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	var b models.Book
	if err := db.DB.First(&b, "id = ?", id).Error; err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	if b.Owner != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var input models.Book
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// обновляем поля (не трогаем owner и id)
	b.Title = input.Title
	b.Author = input.Author
	b.Year = input.Year
	b.Genre = input.Genre

	if err := db.DB.Save(&b).Error; err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(b)
}

// DELETE /books/{id}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	var b models.Book
	if err := db.DB.First(&b, "id = ?", id).Error; err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	if b.Owner != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := db.DB.Delete(&b).Error; err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
