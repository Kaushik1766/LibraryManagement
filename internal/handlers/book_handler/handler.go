package bookhandler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	bookservice "github.com/Kaushik1766/LibraryManagement/internal/service/book_service"
	weberrors "github.com/Kaushik1766/LibraryManagement/internal/web_errors"
)

type BookHandler struct {
	bookService bookservice.BookManager
}

func NewBookHandler(bookService bookservice.BookManager) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (handler *BookHandler) AddBook(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req models.AddBookDTO

	data, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(data, &req)
	if err != nil {
		weberrors.SendError(err, http.StatusBadRequest, w)
		return
	}

	err = handler.bookService.AddBook(ctx, req)
	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *BookHandler) GetAllBooks(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	books, err := handler.bookService.GetAllBooks(ctx, query.Get("title"), query.Get("author"))
	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}
