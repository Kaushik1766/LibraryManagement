package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID
	Book       Book
	User       User
	IssuedAt   time.Time
	IssuedTill time.Time
	ReturnedAt *time.Time
}

type TransactionDTO struct {
	ID         string `json:"transaction_id"`
	BookID     string `json:"book_id"`
	BookName   string `json:"book_name"`
	UserEmail  string `json:"user_email"`
	IssuedAt   string `json:"issued_at"`
	IssuedTill string `json:"issued_till"`
	ReturnedAt string `json:"returned_at"`
}

type GetTransactionRequestDTO struct {
	TransactionId string
	UserId        string
	StartTime     string
	EndTime       string
	Returned      string
	BookName      string
}

type OverdueTransactionDTO struct {
	ID         string `json:"transaction_id"`
	BookID     string `json:"book_id"`
	BookName   string `json:"book_name"`
	IssuedAt   string `json:"issued_at"`
	IssuedTill string `json:"issued_till"`
	ReturnedAt string `json:"returned_at"`
}
