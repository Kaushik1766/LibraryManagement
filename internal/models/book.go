package models

import "github.com/google/uuid"

type Book struct {
	ID       uuid.UUID
	Title    string
	Author   string
	IssuedTo *User
}

type AddBookDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Copies int    `json:"copies"`
}

type BookDTO struct {
	ID       string `json:"book_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	IssuedTo string `json:"issued_to,omitempty"`
}
