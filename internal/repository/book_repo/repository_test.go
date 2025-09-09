package bookrepo

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kaushik1766/LibraryManagement/internal/models"
	"github.com/google/uuid"
)

func TestBookRepository_AddBook(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		title  string
		author string
		copies int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid add book",
			fields: fields{
				db: db,
			},
			args: args{
				title:  "asdf",
				author: "asdf",
				copies: 2,
			},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectExec(`(?i)insert into books.*`).
					WithArgs("asdf", "asdf", 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "invalid add book",
			fields: fields{
				db: db,
			},
			args: args{
				title:  "asdf",
				author: "asdf",
				copies: 2,
			},
			wantErr: true,
			mockSetup: func() {
				mock.ExpectExec(`(?i)insert into books.*`).
					WithArgs("asdf", "asdf", 2).
					WillReturnError(errors.New("invalid add book"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &BookRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			if err := repo.AddBook(tt.args.title, tt.args.author, tt.args.copies); (err != nil) != tt.wantErr {
				t.Errorf("AddBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookRepository_GetAllBooks(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	book1 := models.Book{
		ID:       uuid.New(),
		Title:    "harry potter",
		Author:   "jk rowling",
		IssuedTo: nil,
	}
	book2 := models.Book{
		ID:     uuid.New(),
		Title:  "harry potter",
		Author: "jk rowling",
		IssuedTo: &models.User{
			Email: "kaushik@a.com",
		},
	}

	type fields struct {
		db *sql.DB
	}
	type args struct {
		title  string
		author string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []models.Book
		wantErr   bool
		mockSetup func()
	}{
		{
			name:   "valid get all books",
			fields: fields{db: db},
			args: args{
				title:  book1.Title,
				author: "",
			},
			want: []models.Book{
				book1,
			},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from books .* left join transactions .* left join users").
					WithArgs(book1.Title, "").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "email"}).
						AddRow(book1.ID, book1.Title, book1.Author, nil))
			},
		},
		{
			name:   "invalid get all books",
			fields: fields{db: db},
			args: args{
				title:  book1.Title,
				author: "",
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from books .* left join transactions .* left join users").
					WithArgs(book1.Title, "").
					WillReturnError(errors.New("error retrieving books"))
			},
		},
		{
			name:   "book issued to user",
			fields: fields{db: db},
			args: args{
				title:  book2.Title,
				author: "",
			},
			want: []models.Book{
				book2,
			},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from books .* left join transactions .* left join users").
					WithArgs(book2.Title, "").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "email"}).
						AddRow(book2.ID, book2.Title, book2.Author, book2.IssuedTo.Email))
			},
		},
		{
			name:   "invalid field in books",
			fields: fields{db: db},
			args: args{
				title:  book1.Title,
				author: "",
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from books .* left join transactions .* left join users").
					WithArgs(book1.Title, "").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "email"}).
						AddRow("asdfasd", book1.Title, book1.Author, nil))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &BookRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			got, err := repo.GetAllBooks(tt.args.title, tt.args.author)
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

func TestNewBookRepository(t *testing.T) {

	db, _, _ := sqlmock.New()
	defer db.Close()

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name      string
		args      args
		want      *BookRepository
		setupMock func()
	}{
		{
			name: "valid",
			args: args{
				db: db,
			},
			want: &BookRepository{
				db: db,
			},
			setupMock: func() {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			if got := NewBookRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBookRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
