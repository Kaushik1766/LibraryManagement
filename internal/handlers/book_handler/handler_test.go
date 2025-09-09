package bookhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	bookservice "github.com/Kaushik1766/LibraryManagement/internal/service/book_service"
	"github.com/Kaushik1766/LibraryManagement/mocks"
	"go.uber.org/mock/gomock"
)

func anyToReader(data any) io.Reader {
	dataJsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(dataJsonBytes)
}

func TestNewBookHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mocks.NewMockBookManager(ctrl)

	type args struct {
		bookService bookservice.BookManager
	}
	tests := []struct {
		name string
		args args
		want *BookHandler
	}{
		{
			name: "valid",
			args: args{
				bookService: mockBookService,
			},
			want: &BookHandler{
				bookService: mockBookService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBookHandler(tt.args.bookService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookHandler_AddBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mocks.NewMockBookManager(ctrl)

	type fields struct {
		bookService bookservice.BookManager
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
		mockSetup      func()
	}{
		{
			name: "valid add book",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/books", anyToReader(models.AddBookDTO{
					Title:  "Harry Potter",
					Author: "J.K. Rowling",
					Copies: 5,
				})),
			},
			expectedStatus: http.StatusCreated,
			mockSetup: func() {
				mockBookService.EXPECT().AddBook(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "invalid json",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader([]byte("invalid json"))),
			},
			expectedStatus: http.StatusBadRequest,
			mockSetup: func() {
			},
		},
		{
			name: "service error",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/books", anyToReader(models.AddBookDTO{
					Title:  "Harry Potter",
					Author: "J.K. Rowling",
					Copies: 5,
				})),
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockBookService.EXPECT().AddBook(gomock.Any(), gomock.Any()).Return(errors.New("service error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &BookHandler{
				bookService: tt.fields.bookService,
			}
			tt.mockSetup()
			handler.AddBook(context.Background(), tt.args.w, tt.args.r)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("AddBook() status = %v, want %v", recorder.Code, tt.expectedStatus)
				}
			}
		})
	}
}

func TestBookHandler_GetAllBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mocks.NewMockBookManager(ctrl)

	type fields struct {
		bookService bookservice.BookManager
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus int
		mockSetup      func()
	}{
		{
			name: "valid get all books",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/books", nil),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockBookService.EXPECT().GetAllBooks(gomock.Any(), "", "").Return([]models.BookDTO{
					{
						ID:       "550e8400-e29b-41d4-a716-446655440000",
						Title:    "Harry Potter",
						Author:   "J.K. Rowling",
						IssuedTo: "",
					},
				}, nil)
			},
		},
		{
			name: "valid get books with title filter",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/books?title=Harry%20Potter", nil),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockBookService.EXPECT().GetAllBooks(gomock.Any(), "Harry Potter", "").Return([]models.BookDTO{
					{
						ID:       "550e8400-e29b-41d4-a716-446655440000",
						Title:    "Harry Potter",
						Author:   "J.K. Rowling",
						IssuedTo: "",
					},
				}, nil)
			},
		},
		{
			name: "valid get books with author filter",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/books?author=J.K.%20Rowling", nil),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockBookService.EXPECT().GetAllBooks(gomock.Any(), "", "J.K. Rowling").Return([]models.BookDTO{
					{
						ID:       "550e8400-e29b-41d4-a716-446655440000",
						Title:    "Harry Potter",
						Author:   "J.K. Rowling",
						IssuedTo: "",
					},
				}, nil)
			},
		},
		{
			name: "valid get books with both filters",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/books?title=Harry%20Potter&author=J.K.%20Rowling", nil),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockBookService.EXPECT().GetAllBooks(gomock.Any(), "Harry Potter", "J.K. Rowling").Return([]models.BookDTO{
					{
						ID:       "550e8400-e29b-41d4-a716-446655440000",
						Title:    "Harry Potter",
						Author:   "J.K. Rowling",
						IssuedTo: "",
					},
				}, nil)
			},
		},
		{
			name: "service error",
			fields: fields{
				bookService: mockBookService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/books", nil),
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockBookService.EXPECT().GetAllBooks(gomock.Any(), "", "").Return(nil, errors.New("service error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &BookHandler{
				bookService: tt.fields.bookService,
			}
			tt.mockSetup()
			handler.GetAllBooks(context.Background(), tt.args.w, tt.args.r)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("GetAllBooks() status = %v, want %v", recorder.Code, tt.expectedStatus)
				}
			}
		})
	}
}
