package bookservice

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	"github.com/Kaushik1766/LibraryManagement/internal/models/enums/roles"
	bookrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/book_repo"
	"github.com/Kaushik1766/LibraryManagement/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestBookService_AddBook(t *testing.T) {

	ctrl := gomock.NewController(t)

	mockBookRepo := mocks.NewMockBookStorage(ctrl)

	type fields struct {
		bookRepo bookrepo.BookStorage
	}
	type args struct {
		ctx     context.Context
		bookReq models.AddBookDTO
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name:   "invalid add book",
			fields: fields{mockBookRepo},
			args: args{
				ctx: context.Background(),
				bookReq: models.AddBookDTO{
					Title:  "",
					Author: "",
					Copies: 4,
				},
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name:   "invalid context",
			fields: fields{mockBookRepo},
			args: args{
				ctx: context.Background(),
				bookReq: models.AddBookDTO{
					Title:  "asdfa",
					Author: "adfsadf",
					Copies: 4,
				},
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name:   "unauthorised user",
			fields: fields{mockBookRepo},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Role: roles.Customer,
				}),
				bookReq: models.AddBookDTO{
					Title:  "asdfa",
					Author: "adfsadf",
					Copies: 4,
				},
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name:   "authorised user",
			fields: fields{mockBookRepo},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Role: roles.Staff,
				}),
				bookReq: models.AddBookDTO{
					Title:  "asdfa",
					Author: "adfsadf",
					Copies: 4,
				},
			},
			wantErr: false,
			mockSetup: func() {
				mockBookRepo.EXPECT().AddBook("asdfa", "adfsadf", 4).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &BookService{
				bookRepo: tt.fields.bookRepo,
			}
			tt.mockSetup()
			if err := service.AddBook(tt.args.ctx, tt.args.bookReq); (err != nil) != tt.wantErr {
				t.Errorf("AddBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookService_GetAllBooks(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockBookRepo := mocks.NewMockBookStorage(ctrl)

	book1 := models.Book{
		ID:       uuid.New(),
		Title:    "asdf",
		Author:   "asdf",
		IssuedTo: nil,
	}

	book2 := models.Book{
		ID:     uuid.New(),
		Title:  "asdf",
		Author: "asdf",
		IssuedTo: &models.User{
			Email: "kaushik@a.com",
		},
	}
	type fields struct {
		bookRepo bookrepo.BookStorage
	}
	type args struct {
		ctx    context.Context
		title  string
		author string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []models.BookDTO
		wantErr   bool
		setupMock func()
	}{
		{
			name:   "invalid context",
			fields: fields{bookRepo: mockBookRepo},
			args: args{
				ctx:    context.Background(),
				title:  "asdf",
				author: "asdf",
			},
			want:    nil,
			wantErr: true,
			setupMock: func() {
			},
		},
		{
			name:   "repo error",
			fields: fields{bookRepo: mockBookRepo},
			args: args{
				ctx:    context.WithValue(context.Background(), "user", models.UserJwt{Role: roles.Customer}),
				title:  "asdf",
				author: "asdf",
			},
			want:    nil,
			wantErr: true,
			setupMock: func() {
				mockBookRepo.EXPECT().GetAllBooks(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
		},
		{
			name:   "customer get all books",
			fields: fields{bookRepo: mockBookRepo},
			args: args{
				ctx:    context.WithValue(context.Background(), "user", models.UserJwt{Role: roles.Customer}),
				title:  "asdf",
				author: "asdf",
			},
			want: []models.BookDTO{
				{
					ID:     book1.ID.String(),
					Title:  book1.Title,
					Author: book1.Author,
				},
			},
			wantErr: false,
			setupMock: func() {
				mockBookRepo.EXPECT().GetAllBooks(gomock.Any(), gomock.Any()).Return([]models.Book{
					book1,
				}, nil)
			},
		},
		{
			name:   "staff get all books",
			fields: fields{bookRepo: mockBookRepo},
			args: args{
				ctx:    context.WithValue(context.Background(), "user", models.UserJwt{Role: roles.Staff}),
				title:  "asdf",
				author: "asdf",
			},
			want: []models.BookDTO{
				{
					ID:       book1.ID.String(),
					Title:    book1.Title,
					Author:   book1.Author,
					IssuedTo: "none",
				},
				{
					ID:       book2.ID.String(),
					Title:    book2.Title,
					Author:   book2.Author,
					IssuedTo: book2.IssuedTo.Email,
				},
			},
			wantErr: false,
			setupMock: func() {
				mockBookRepo.EXPECT().GetAllBooks(gomock.Any(), gomock.Any()).Return([]models.Book{
					book1,
					book2,
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &BookService{
				bookRepo: tt.fields.bookRepo,
			}
			tt.setupMock()
			got, err := service.GetAllBooks(tt.args.ctx, tt.args.title, tt.args.author)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllBooks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBookService(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockBookRepo := mocks.NewMockBookStorage(ctrl)
	type args struct {
		bookRepo bookrepo.BookStorage
	}
	tests := []struct {
		name string
		args args
		want *BookService
	}{
		{
			name: "valid",
			args: args{mockBookRepo},
			want: &BookService{bookRepo: mockBookRepo},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBookService(tt.args.bookRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookService() = %v, want %v", got, tt.want)
			}
		})
	}
}
