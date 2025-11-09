package models

type Book struct {
	ID     string `gorm:"primaryKey;type:uuid" json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
	Owner  uint   `json:"owner"`
}

type CreateBookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	Year   int    `json:"year" validate:"gte=0,lte=2100"`
	Genre  string `json:"genre" validate:"required"`
}

type UpdateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
}

type BookResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
}
