package bookservice

import (
	"context"
	"errors"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	"github.com/Kaushik1766/LibraryManagement/internal/models/enums/roles"
	bookrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/book_repo"
)

type BookService struct {
	bookRepo bookrepo.BookStorage
}

func NewBookService(bookRepo bookrepo.BookStorage) *BookService {
	return &BookService{
		bookRepo: bookRepo,
	}
}

func (service *BookService) AddBook(ctx context.Context, bookReq models.AddBookDTO) error {
	if bookReq.Title == "" || bookReq.Author == "" || bookReq.Copies <= 0 {
		return errors.New("invalid input")
	}

	userCtx, ok := ctx.Value("user").(models.UserJwt)
	if !ok {
		return errors.New("invalid context")
	}

	if userCtx.Role != roles.Staff {
		return errors.New("unauthorised user")
	}

	return service.bookRepo.AddBook(bookReq.Title, bookReq.Author, bookReq.Copies)
}

func (service *BookService) GetAllBooks(ctx context.Context, title, author string) ([]models.BookDTO, error) {
	userCtx, ok := ctx.Value("user").(models.UserJwt)
	if !ok {
		return nil, errors.New("invalid context")
	}

	books, err := service.bookRepo.GetAllBooks(title, author)
	if err != nil {
		return nil, err
	}

	var bookResponse []models.BookDTO
	if userCtx.Role == roles.Staff {
		for _, val := range books {
			if val.IssuedTo == nil {
				bookResponse = append(bookResponse, models.BookDTO{
					ID:       val.ID.String(),
					Title:    val.Title,
					Author:   val.Author,
					IssuedTo: "none",
				})
			} else {
				bookResponse = append(bookResponse, models.BookDTO{
					ID:       val.ID.String(),
					Title:    val.Title,
					Author:   val.Author,
					IssuedTo: val.IssuedTo.Email,
				})
			}
		}
	} else {
		for _, val := range books {
			if val.IssuedTo == nil {
				bookResponse = append(bookResponse, models.BookDTO{
					ID:     val.ID.String(),
					Title:  val.Title,
					Author: val.Author,
				})
			}
		}
	}

	return bookResponse, nil
}
