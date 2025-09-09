package bookservice

import (
	"context"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/mock_book_manager.go -package=mocks
type BookManager interface {
	AddBook(ctx context.Context, bookReq models.AddBookDTO) error
	GetAllBooks(ctx context.Context, title, author string) ([]models.BookDTO, error)
}
