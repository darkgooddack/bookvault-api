package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/darkgooddack/bookvault-api/db"
	"github.com/darkgooddack/bookvault-api/middleware"
	"github.com/darkgooddack/bookvault-api/models"
	"github.com/go-playground/validator/v10"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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

	var resp []models.BookResponse
	for _, b := range books {
		resp = append(resp, models.BookResponse{
			ID:     b.ID,
			Title:  b.Title,
			Author: b.Author,
			Year:   b.Year,
			Genre:  b.Genre,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r)

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book := models.Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
		Genre:  req.Genre,
		Owner:  userID,
	}

	if err := db.DB.Create(&book).Error; err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

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

	resp := models.BookResponse{
		ID:     b.ID,
		Title:  b.Title,
		Author: b.Author,
		Year:   b.Year,
		Genre:  b.Genre,
	}

	json.NewEncoder(w).Encode(resp)
}

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

	var req models.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// обновляем поля
	if req.Title != "" {
		b.Title = req.Title
	}
	if req.Author != "" {
		b.Author = req.Author
	}
	if req.Year != 0 {
		b.Year = req.Year
	}
	if req.Genre != "" {
		b.Genre = req.Genre
	}

	if err := db.DB.Save(&b).Error; err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	resp := models.BookResponse{
		ID:     b.ID,
		Title:  b.Title,
		Author: b.Author,
		Year:   b.Year,
		Genre:  b.Genre,
	}

	json.NewEncoder(w).Encode(resp)
}

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
