package transactionrepo

import "github.com/Kaushik1766/LibraryManagement/internal/models"

//go:generate mockgen -source=interface.go -destination=../../../mocks/mock_transaction_storage.go -package=mocks
type TransactionStorage interface {
	IssueBook(bookId, userId, issueFor string) (string, error)
	ReturnBook(bookId, userId string) error
	GetAllTransactions(dto models.GetTransactionRequestDTO) ([]models.Transaction, error)
	GetOverDueTransactions(userId string) ([]models.Transaction, error)
}
