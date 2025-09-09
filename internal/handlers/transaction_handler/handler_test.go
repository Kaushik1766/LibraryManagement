package transactionhandler

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
	transactionservice "github.com/Kaushik1766/LibraryManagement/internal/service/transaction_service"
	"github.com/Kaushik1766/LibraryManagement/mocks"
	"go.uber.org/mock/gomock"
)

func anyToReader(data any) io.Reader {
	dataJsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(dataJsonBytes)
}

func TestNewTransactionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mocks.NewMockTransactionManager(ctrl)

	type args struct {
		transactionService transactionservice.TransactionManager
	}
	tests := []struct {
		name string
		args args
		want *TransactionHandler
	}{
		{
			name: "valid",
			args: args{
				transactionService: mockTransactionService,
			},
			want: &TransactionHandler{
				transactionService: mockTransactionService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionHandler(tt.args.transactionService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionHandler_IssueBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mocks.NewMockTransactionManager(ctrl)

	type fields struct {
		transactionService transactionservice.TransactionManager
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
			name: "valid issue book",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transactions/issue", anyToReader(map[string]string{
					"book_id":   "550e8400-e29b-41d4-a716-446655440000",
					"issue_for": "7 days",
				})),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockTransactionService.EXPECT().IssueBook(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000", "7 days").Return("550e8400-e29b-41d4-a716-446655440001", nil)
			},
		},
		{
			name: "invalid json",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transactions/issue", bytes.NewReader([]byte("invalid json"))),
			},
			expectedStatus: http.StatusBadRequest,
			mockSetup: func() {
			},
		},
		{
			name: "service error",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transactions/issue", anyToReader(map[string]string{
					"book_id":   "550e8400-e29b-41d4-a716-446655440000",
					"issue_for": "7 days",
				})),
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockTransactionService.EXPECT().IssueBook(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000", "7 days").Return("", errors.New("service error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionHandler{
				transactionService: tt.fields.transactionService,
			}
			tt.mockSetup()
			handler.IssueBook(context.Background(), tt.args.w, tt.args.r)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("IssueBook() status = %v, want %v", recorder.Code, tt.expectedStatus)
				}
			}
		})
	}
}

func TestTransactionHandler_ReturnBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mocks.NewMockTransactionManager(ctrl)

	type fields struct {
		transactionService transactionservice.TransactionManager
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
			name: "valid return book",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transactions/return", anyToReader(map[string]string{
					"book_id": "550e8400-e29b-41d4-a716-446655440000",
				})),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockTransactionService.EXPECT().ReturnBook(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").Return(nil)
			},
		},
		{
			name: "invalid json",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transactions/return", bytes.NewReader([]byte("invalid json"))),
			},
			expectedStatus: http.StatusBadRequest,
			mockSetup: func() {
			},
		},
		{
			name: "service error",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/transactions/return", anyToReader(map[string]string{
					"book_id": "550e8400-e29b-41d4-a716-446655440000",
				})),
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockTransactionService.EXPECT().ReturnBook(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").Return(errors.New("service error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionHandler{
				transactionService: tt.fields.transactionService,
			}
			tt.mockSetup()
			handler.ReturnBook(context.Background(), tt.args.w, tt.args.r)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("ReturnBook() status = %v, want %v", recorder.Code, tt.expectedStatus)
				}
			}
		})
	}
}

func TestTransactionHandler_GetAllTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mocks.NewMockTransactionManager(ctrl)

	type fields struct {
		transactionService transactionservice.TransactionManager
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
			name: "valid get all transactions",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/transactions", nil),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockTransactionService.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).Return([]models.TransactionDTO{
					{
						ID:         "550e8400-e29b-41d4-a716-446655440001",
						BookID:     "550e8400-e29b-41d4-a716-446655440000",
						BookName:   "Harry Potter",
						UserEmail:  "kaushik@a.com",
						IssuedAt:   "2025-09-04 03:00:43 +0530 IST",
						IssuedTill: "2025-09-11 03:00:43 +0530 IST",
						ReturnedAt: "",
					},
				}, nil)
			},
		},
		{
			name: "valid get transactions with filters",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/transactions?startTime=2025-09-01T00:00:00Z&endTime=2025-09-30T23:59:59Z&title=Harry%20Potter", nil),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockTransactionService.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).Return([]models.TransactionDTO{
					{
						ID:         "550e8400-e29b-41d4-a716-446655440001",
						BookID:     "550e8400-e29b-41d4-a716-446655440000",
						BookName:   "Harry Potter",
						UserEmail:  "kaushik@a.com",
						IssuedAt:   "2025-09-04 03:00:43 +0530 IST",
						IssuedTill: "2025-09-11 03:00:43 +0530 IST",
						ReturnedAt: "",
					},
				}, nil)
			},
		},
		{
			name: "service error",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/transactions", nil),
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockTransactionService.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionHandler{
				transactionService: tt.fields.transactionService,
			}
			tt.mockSetup()
			handler.GetAllTransactions(context.Background(), tt.args.w, tt.args.r)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("GetAllTransactions() status = %v, want %v", recorder.Code, tt.expectedStatus)
				}
			}
		})
	}
}

func TestTransactionHandler_GetOverdueTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mocks.NewMockTransactionManager(ctrl)

	type fields struct {
		transactionService transactionservice.TransactionManager
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
			name: "valid get overdue transactions",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/transactions/overdue", nil),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockTransactionService.EXPECT().GetOverdueTransactions(gomock.Any()).Return([]models.OverdueTransactionDTO{
					{
						ID:         "550e8400-e29b-41d4-a716-446655440001",
						BookID:     "550e8400-e29b-41d4-a716-446655440000",
						BookName:   "Harry Potter",
						IssuedAt:   "2025-08-30 03:00:43 +0530 IST",
						IssuedTill: "2025-09-04 03:00:43 +0530 IST",
						ReturnedAt: "not yet returned",
					},
				}, nil)
			},
		},
		{
			name: "service error",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/transactions/overdue", nil),
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockTransactionService.EXPECT().GetOverdueTransactions(gomock.Any()).Return(nil, errors.New("service error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionHandler{
				transactionService: tt.fields.transactionService,
			}
			tt.mockSetup()
			handler.GetOverdueTransactions(context.Background(), tt.args.w, tt.args.r)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("GetOverdueTransactions() status = %v, want %v", recorder.Code, tt.expectedStatus)
				}
			}
		})
	}
}

func TestTransactionHandler_GetTransactionById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mocks.NewMockTransactionManager(ctrl)

	type fields struct {
		transactionService transactionservice.TransactionManager
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
			name: "valid get transaction by id",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/transactions/550e8400-e29b-41d4-a716-446655440001", nil)
					req.SetPathValue("transactionId", "550e8400-e29b-41d4-a716-446655440001")
					return req
				}(),
			},
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				mockTransactionService.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).Return([]models.TransactionDTO{
					{
						ID:         "550e8400-e29b-41d4-a716-446655440001",
						BookID:     "550e8400-e29b-41d4-a716-446655440000",
						BookName:   "Harry Potter",
						UserEmail:  "kaushik@a.com",
						IssuedAt:   "2025-09-04 03:00:43 +0530 IST",
						IssuedTill: "2025-09-11 03:00:43 +0530 IST",
						ReturnedAt: "",
					},
				}, nil)
			},
		},
		{
			name: "service error",
			fields: fields{
				transactionService: mockTransactionService,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/transactions/550e8400-e29b-41d4-a716-446655440001", nil)
					req.SetPathValue("transactionId", "550e8400-e29b-41d4-a716-446655440001")
					return req
				}(),
			},
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func() {
				mockTransactionService.EXPECT().GetTransactions(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &TransactionHandler{
				transactionService: tt.fields.transactionService,
			}
			tt.mockSetup()
			handler.GetTransactionById(context.Background(), tt.args.w, tt.args.r)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("GetTransactionById() status = %v, want %v", recorder.Code, tt.expectedStatus)
				}
			}
		})
	}
}
