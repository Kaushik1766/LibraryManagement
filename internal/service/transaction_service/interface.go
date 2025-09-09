package transactionservice

import (
	"context"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/mock_transaction_manager.go -package=mocks
type TransactionManager interface {
	IssueBook(ctx context.Context, bookId, issueFor string) (string, error)
	ReturnBook(ctx context.Context, bookId string) error
	GetTransactions(ctx context.Context, dto models.GetTransactionRequestDTO) ([]models.TransactionDTO, error)
	GetOverdueTransactions(ctx context.Context) ([]models.OverdueTransactionDTO, error)
}
