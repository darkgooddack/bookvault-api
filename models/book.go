package models

type Book struct {
	ID     string `gorm:"primaryKey;type:uuid" json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
	Owner  uint   `json:"owner"`
}
