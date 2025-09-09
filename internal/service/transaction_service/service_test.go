package transactionservice

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	"github.com/Kaushik1766/LibraryManagement/internal/models/enums/roles"
	bookrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/book_repo"
	transactionrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/transaction_repo"
	"github.com/Kaushik1766/LibraryManagement/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestTransactionService_GetOverdueTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookRepo := mocks.NewMockBookStorage(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionStorage(ctrl)

	type fields struct {
		bookRepo        bookrepo.BookStorage
		transactionRepo transactionrepo.TransactionStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []models.OverdueTransactionDTO
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid staff overdue transactions",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "staff@example.com",
					Role:  roles.Staff,
				}),
			},
			want: []models.OverdueTransactionDTO{
				{
					ID:         "550e8400-e29b-41d4-a716-446655440000",
					BookID:     "550e8400-e29b-41d4-a716-446655440001",
					BookName:   "Test Book",
					IssuedAt:   "2025-08-30 03:00:43 +0530 IST",
					IssuedTill: "2025-09-04 03:00:43 +0530 IST",
					ReturnedAt: "not yet returned",
				},
			},
			wantErr: false,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().GetOverDueTransactions("").Return([]models.Transaction{
					{
						ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						Book:       models.Book{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"), Title: "Test Book"},
						User:       models.User{Email: "user@example.com"},
						IssuedAt:   time.Date(2025, 8, 30, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						IssuedTill: time.Date(2025, 9, 4, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						ReturnedAt: nil,
					},
				}, nil)
			},
		},
		{
			name: "valid customer overdue transactions",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
			},
			want: []models.OverdueTransactionDTO{
				{
					ID:         "550e8400-e29b-41d4-a716-446655440002",
					BookID:     "550e8400-e29b-41d4-a716-446655440003",
					BookName:   "Test Book",
					IssuedAt:   "2025-08-30 03:00:43 +0530 IST",
					IssuedTill: "2025-09-04 03:00:43 +0530 IST",
					ReturnedAt: "not yet returned",
				},
			},
			wantErr: false,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().GetOverDueTransactions(gomock.Any()).Return([]models.Transaction{
					{
						ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
						Book:       models.Book{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"), Title: "Test Book"},
						User:       models.User{Email: "customer@example.com"},
						IssuedAt:   time.Date(2025, 8, 30, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						IssuedTill: time.Date(2025, 9, 4, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						ReturnedAt: nil,
					},
				}, nil)
			},
		},
		{
			name: "invalid user context",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "repository error",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "staff@example.com",
					Role:  roles.Staff,
				}),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().GetOverDueTransactions("").Return(nil, errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TransactionService{
				bookRepo:        tt.fields.bookRepo,
				transactionRepo: tt.fields.transactionRepo,
			}
			tt.mockSetup()
			got, err := service.GetOverdueTransactions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionService.GetOverdueTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionService.GetOverdueTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTransactionService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookRepo := mocks.NewMockBookStorage(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionStorage(ctrl)

	type args struct {
		bookRepo        bookrepo.BookStorage
		transactionRepo transactionrepo.TransactionStorage
	}
	tests := []struct {
		name string
		args args
		want *TransactionService
	}{
		{
			name: "valid",
			args: args{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			want: &TransactionService{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionService(tt.args.bookRepo, tt.args.transactionRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_IssueBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookRepo := mocks.NewMockBookStorage(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionStorage(ctrl)

	type fields struct {
		bookRepo        bookrepo.BookStorage
		transactionRepo transactionrepo.TransactionStorage
	}
	type args struct {
		ctx      context.Context
		bookId   string
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
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				bookId:   uuid.New().String(),
				issueFor: "7 days",
			},
			want:    "550e8400-e29b-41d4-a716-446655440004",
			wantErr: false,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().IssueBook(gomock.Any(), gomock.Any(), "7 days").Return("550e8400-e29b-41d4-a716-446655440004", nil)
			},
		},
		{
			name: "valid issue book with empty issueFor",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				bookId:   uuid.New().String(),
				issueFor: "",
			},
			want:    "550e8400-e29b-41d4-a716-446655440005",
			wantErr: false,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().IssueBook(gomock.Any(), gomock.Any(), "1 day").Return("550e8400-e29b-41d4-a716-446655440005", nil)
			},
		},
		{
			name: "invalid user context",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx:      context.Background(),
				bookId:   uuid.New().String(),
				issueFor: "7 days",
			},
			want:    "",
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "staff cannot issue book",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "staff@example.com",
					Role:  roles.Staff,
				}),
				bookId:   uuid.New().String(),
				issueFor: "7 days",
			},
			want:    "",
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "invalid book id",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				bookId:   "",
				issueFor: "7 days",
			},
			want:    "",
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "repository error",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				bookId:   uuid.New().String(),
				issueFor: "7 days",
			},
			want:    "",
			wantErr: true,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().IssueBook(gomock.Any(), gomock.Any(), "7 days").Return("", errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TransactionService{
				bookRepo:        tt.fields.bookRepo,
				transactionRepo: tt.fields.transactionRepo,
			}
			tt.mockSetup()
			got, err := service.IssueBook(tt.args.ctx, tt.args.bookId, tt.args.issueFor)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionService.IssueBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TransactionService.IssueBook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_ReturnBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookRepo := mocks.NewMockBookStorage(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionStorage(ctrl)

	type fields struct {
		bookRepo        bookrepo.BookStorage
		transactionRepo transactionrepo.TransactionStorage
	}
	type args struct {
		ctx    context.Context
		bookId string
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
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				bookId: uuid.New().String(),
			},
			wantErr: false,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().ReturnBook(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "invalid user context",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx:    context.Background(),
				bookId: uuid.New().String(),
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "staff cannot return book",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "staff@example.com",
					Role:  roles.Staff,
				}),
				bookId: uuid.New().String(),
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "invalid book id",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				bookId: "",
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "repository error",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				bookId: uuid.New().String(),
			},
			wantErr: true,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().ReturnBook(gomock.Any(), gomock.Any()).Return(errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TransactionService{
				bookRepo:        tt.fields.bookRepo,
				transactionRepo: tt.fields.transactionRepo,
			}
			tt.mockSetup()
			if err := service.ReturnBook(tt.args.ctx, tt.args.bookId); (err != nil) != tt.wantErr {
				t.Errorf("TransactionService.ReturnBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionService_GetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookRepo := mocks.NewMockBookStorage(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionStorage(ctrl)

	type fields struct {
		bookRepo        bookrepo.BookStorage
		transactionRepo transactionrepo.TransactionStorage
	}
	type args struct {
		ctx context.Context
		dto models.GetTransactionRequestDTO
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []models.TransactionDTO
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid get transactions",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				dto: models.GetTransactionRequestDTO{
					StartTime: time.Now().AddDate(0, -1, 0).Format(time.RFC3339),
					EndTime:   time.Now().Format(time.RFC3339),
				},
			},
			want: []models.TransactionDTO{
				{
					ID:         "550e8400-e29b-41d4-a716-446655440006",
					BookID:     "550e8400-e29b-41d4-a716-446655440007",
					BookName:   "Test Book",
					UserEmail:  "customer@example.com",
					IssuedAt:   "2025-09-04 03:00:43 +0530 IST",
					IssuedTill: "2025-09-11 03:00:43 +0530 IST",
					ReturnedAt: "not returned yet",
				},
			},
			wantErr: false,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().GetAllTransactions(gomock.Any()).Return([]models.Transaction{
					{
						ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
						Book:       models.Book{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"), Title: "Test Book"},
						User:       models.User{Email: "customer@example.com"},
						IssuedAt:   time.Date(2025, 9, 4, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						IssuedTill: time.Date(2025, 9, 11, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						ReturnedAt: nil,
					},
				}, nil)
			},
		},
		{
			name: "valid get transactions with empty times",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				dto: models.GetTransactionRequestDTO{},
			},
			want: []models.TransactionDTO{
				{
					ID:         "550e8400-e29b-41d4-a716-446655440008",
					BookID:     "550e8400-e29b-41d4-a716-446655440009",
					BookName:   "Test Book",
					UserEmail:  "customer@example.com",
					IssuedAt:   "2025-09-04 03:00:43 +0530 IST",
					IssuedTill: "2025-09-11 03:00:43 +0530 IST",
					ReturnedAt: "2025-09-10 03:00:43 +0530 IST",
				},
			},
			wantErr: false,
			mockSetup: func() {
				returnedAt := time.Date(2025, 9, 10, 3, 0, 43, 0, time.FixedZone("IST", 19800))
				mockTransactionRepo.EXPECT().GetAllTransactions(gomock.Any()).Return([]models.Transaction{
					{
						ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440008"),
						Book:       models.Book{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440009"), Title: "Test Book"},
						User:       models.User{Email: "customer@example.com"},
						IssuedAt:   time.Date(2025, 9, 4, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						IssuedTill: time.Date(2025, 9, 11, 3, 0, 43, 0, time.FixedZone("IST", 19800)),
						ReturnedAt: &returnedAt,
					},
				}, nil)
			},
		},
		{
			name: "invalid user context",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.Background(),
				dto: models.GetTransactionRequestDTO{},
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name: "repository error",
			fields: fields{
				bookRepo:        mockBookRepo,
				transactionRepo: mockTransactionRepo,
			},
			args: args{
				ctx: context.WithValue(context.Background(), "user", models.UserJwt{
					Email: "customer@example.com",
					Role:  roles.Customer,
				}),
				dto: models.GetTransactionRequestDTO{},
			},
			want:    nil,
			wantErr: true,
			mockSetup: func() {
				mockTransactionRepo.EXPECT().GetAllTransactions(gomock.Any()).Return(nil, errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TransactionService{
				bookRepo:        tt.fields.bookRepo,
				transactionRepo: tt.fields.transactionRepo,
			}
			tt.mockSetup()
			got, err := service.GetTransactions(tt.args.ctx, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionService.GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionService.GetTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}
