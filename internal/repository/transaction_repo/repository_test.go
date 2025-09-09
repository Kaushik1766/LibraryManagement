package transactionrepo

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kaushik1766/LibraryManagement/internal/models"
	"github.com/google/uuid"
)

func TestNewTransactionRepository(t *testing.T) {

	db, _, _ := sqlmock.New()
	defer db.Close()

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name      string
		args      args
		want      *TransactionRepository
		mockSetup func()
	}{
		{
			name: "valid",
			args: args{
				db: db,
			},
			want: &TransactionRepository{db: db},
			mockSetup: func() {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			if got := NewTransactionRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepository_GetAllTransactions(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	transaction1 := models.Transaction{
		ID:         uuid.New(),
		Book:       models.Book{ID: uuid.New(), Title: "harry potter"},
		User:       models.User{Email: "kaushik@a.com"},
		IssuedAt:   time.Now(),
		IssuedTill: time.Now().AddDate(0, 0, 1),
		ReturnedAt: nil,
	}

	now := time.Now()
	transaction2 := models.Transaction{
		ID:         uuid.New(),
		Book:       models.Book{ID: uuid.New(), Title: "harry potter"},
		User:       models.User{Email: "kaushik@a.com"},
		IssuedAt:   time.Now(),
		IssuedTill: time.Now().AddDate(0, 0, 1),
		ReturnedAt: &now,
	}
	type fields struct {
		db *sql.DB
	}
	type args struct {
		dto models.GetTransactionRequestDTO
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []models.Transaction
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid get transactions",
			fields: fields{
				db: db,
			},
			args: args{
				dto: models.GetTransactionRequestDTO{
					TransactionId: transaction1.ID.String(),
					UserId:        "",
					StartTime:     "",
					EndTime:       "",
					Returned:      "",
					BookName:      "",
				},
			},
			want:    []models.Transaction{transaction1},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join users .* left join books .*").WillReturnRows(sqlmock.NewRows([]string{"id", "issuedAt", "returnedAt", "issuedTill", "bookId", "title", "email"}).AddRow(transaction1.ID, transaction1.IssuedAt, transaction1.ReturnedAt, transaction1.IssuedTill, transaction1.Book.ID, transaction1.Book.Title, transaction1.User.Email))
			},
		},
		{
			name: "invalid get transactions",
			fields: fields{
				db: db,
			},
			args: args{
				dto: models.GetTransactionRequestDTO{
					TransactionId: transaction1.ID.String(),
					UserId:        "",
					StartTime:     "",
					EndTime:       "",
					Returned:      "",
					BookName:      "",
				},
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join users .* left join books .*").WillReturnError(errors.New("invalid query"))
			},
		},
		{
			name: "invalid fields ",
			fields: fields{
				db: db,
			},
			args: args{
				dto: models.GetTransactionRequestDTO{
					TransactionId: transaction1.ID.String(),
					UserId:        "",
					StartTime:     "",
					EndTime:       "",
					Returned:      "",
					BookName:      "",
				},
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join users .* left join books .*").WillReturnRows(sqlmock.NewRows([]string{"id", "issuedAt", "returnedAt", "issuedTill", "bookId", "title", "email"}).AddRow("dd", transaction2.IssuedAt, transaction2.ReturnedAt, transaction2.IssuedTill, transaction2.Book.ID, transaction2.Book.Title, transaction2.User.Email))
			},
		},
		{
			name: "returned at not nil",
			fields: fields{
				db: db,
			},
			args: args{
				dto: models.GetTransactionRequestDTO{
					TransactionId: transaction1.ID.String(),
					UserId:        "",
					StartTime:     "",
					EndTime:       "",
					Returned:      "",
					BookName:      "",
				},
			},
			want:    []models.Transaction{transaction2},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join users .* left join books .*").WillReturnRows(sqlmock.NewRows([]string{"id", "issuedAt", "returnedAt", "issuedTill", "bookId", "title", "email"}).AddRow(transaction2.ID, transaction2.IssuedAt, transaction2.ReturnedAt, transaction2.IssuedTill, transaction2.Book.ID, transaction2.Book.Title, transaction2.User.Email))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TransactionRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			got, err := repo.GetAllTransactions(tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTransactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepository_GetOverDueTransactions(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	transaction1 := models.Transaction{
		ID:         uuid.New(),
		Book:       models.Book{ID: uuid.New(), Title: "harry potter"},
		User:       models.User{},
		IssuedAt:   time.Now(),
		IssuedTill: time.Now().AddDate(0, 0, -1), // Past date for overdue
		ReturnedAt: nil,
	}

	now := time.Now()
	transaction2 := models.Transaction{
		ID:         uuid.New(),
		Book:       models.Book{ID: uuid.New(), Title: "harry potter"},
		User:       models.User{},
		IssuedAt:   time.Now(),
		IssuedTill: time.Now().AddDate(0, 0, -1), // Past date
		ReturnedAt: &now,
	}
	type fields struct {
		db *sql.DB
	}
	type args struct {
		userId string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []models.Transaction
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid get overdue transactions",
			fields: fields{
				db: db,
			},
			args: args{
				userId: uuid.New().String(),
			},
			want:    []models.Transaction{transaction1},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join books .*").WillReturnRows(sqlmock.NewRows([]string{"id", "bookId", "title", "issued_at", "issued_till", "returned_at"}).AddRow(transaction1.ID, transaction1.Book.ID, transaction1.Book.Title, transaction1.IssuedAt, transaction1.IssuedTill, transaction1.ReturnedAt))
			},
		},
		{
			name: "invalid get overdue transactions",
			fields: fields{
				db: db,
			},
			args: args{
				userId: uuid.New().String(),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join books .*").WillReturnError(errors.New("invalid query"))
			},
		},
		{
			name: "returned at not nil",
			fields: fields{
				db: db,
			},
			args: args{
				userId: uuid.New().String(),
			},
			want:    []models.Transaction{transaction2},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join books .*").WillReturnRows(sqlmock.NewRows([]string{"id", "bookId", "title", "issued_at", "issued_till", "returned_at"}).AddRow(transaction2.ID, transaction2.Book.ID, transaction2.Book.Title, transaction2.IssuedAt, transaction2.IssuedTill, transaction2.ReturnedAt))
			},
		},
		{
			name: "empty userId",
			fields: fields{
				db: db,
			},
			args: args{
				userId: "",
			},
			want:    []models.Transaction{transaction1},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join books .*").WillReturnRows(sqlmock.NewRows([]string{"id", "bookId", "title", "issued_at", "issued_till", "returned_at"}).AddRow(transaction1.ID, transaction1.Book.ID, transaction1.Book.Title, transaction1.IssuedAt, transaction1.IssuedTill, transaction1.ReturnedAt))
			},
		},
		{
			name: "scan error",
			fields: fields{
				db: db,
			},
			args: args{
				userId: uuid.New().String(),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from transactions .* left join books .*").WillReturnRows(sqlmock.NewRows([]string{"id", "bookId", "title", "issued_at", "issued_till", "returned_at"}).AddRow("invalid-uuid", transaction1.Book.ID, transaction1.Book.Title, transaction1.IssuedAt, transaction1.IssuedTill, transaction1.ReturnedAt))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TransactionRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			got, err := repo.GetOverDueTransactions(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOverDueTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOverDueTransactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepository_IssueBook(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	transactionId := uuid.New().String()
	bookId := uuid.New().String()
	userId := uuid.New().String()
	issueFor := "7 days"

	type fields struct {
		db *sql.DB
	}
	type args struct {
		bookId   string
		userId   string
		issueFor string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      string
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid issue book",
			fields: fields{
				db: db,
			},
			args: args{
				bookId:   bookId,
				userId:   userId,
				issueFor: issueFor,
			},
			want:    transactionId,
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)insert into transactions .*").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(transactionId))
			},
		},
		{
			name: "book already issued",
			fields: fields{
				db: db,
			},
			args: args{
				bookId:   bookId,
				userId:   userId,
				issueFor: issueFor,
			},
			want:    "",
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)insert into transactions .*").WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "database error",
			fields: fields{
				db: db,
			},
			args: args{
				bookId:   bookId,
				userId:   userId,
				issueFor: issueFor,
			},
			want:    "",
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)insert into transactions .*").WillReturnError(errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TransactionRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			got, err := repo.IssueBook(tt.args.bookId, tt.args.userId, tt.args.issueFor)
			if (err != nil) != tt.wantErr {
				t.Errorf("IssueBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IssueBook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepository_ReturnBook(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	bookId := uuid.New().String()
	userId := uuid.New().String()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		bookId string
		userId string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid return book",
			fields: fields{
				db: db,
			},
			args: args{
				bookId: bookId,
				userId: userId,
			},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectExec("(?i)update transactions .*").WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "book already returned",
			fields: fields{
				db: db,
			},
			args: args{
				bookId: bookId,
				userId: userId,
			},
			wantErr: true,
			mockSetup: func() {
				mock.ExpectExec("(?i)update transactions .*").WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
		{
			name: "database error",
			fields: fields{
				db: db,
			},
			args: args{
				bookId: bookId,
				userId: userId,
			},
			wantErr: true,
			mockSetup: func() {
				mock.ExpectExec("(?i)update transactions .*").WillReturnError(errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TransactionRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			if err := repo.ReturnBook(tt.args.bookId, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("ReturnBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
