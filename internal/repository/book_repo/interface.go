package bookrepo

import "github.com/Kaushik1766/LibraryManagement/internal/models"

//go:generate mockgen -source=interface.go -destination=../../../mocks/mock_book_storage.go -package=mocks
type BookStorage interface {
	AddBook(title, author string, copies int) error
	GetAllBooks(title, author string) ([]models.Book, error)
}
