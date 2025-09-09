package transactionservice

import (
	"context"
	"errors"
	"time"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	"github.com/Kaushik1766/LibraryManagement/internal/models/enums/roles"
	bookrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/book_repo"
	transactionrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/transaction_repo"
)

type TransactionService struct {
	bookRepo        bookrepo.BookStorage
	transactionRepo transactionrepo.TransactionStorage
}

func (service *TransactionService) GetOverdueTransactions(ctx context.Context) ([]models.OverdueTransactionDTO, error) {
	userCtx, ok := ctx.Value("user").(models.UserJwt)

	if !ok {
		return nil, errors.New("invalid user")
	}

	var overdueTransactions []models.Transaction
	var err error
	if userCtx.Role == roles.Staff {
		// staff can get overdue transactions of all users
		overdueTransactions, err = service.transactionRepo.GetOverDueTransactions("")
	} else {
		overdueTransactions, err = service.transactionRepo.GetOverDueTransactions(userCtx.ID)
	}
	if err != nil {
		return nil, err
	}

	var overdueDto []models.OverdueTransactionDTO
	for _, val := range overdueTransactions {
		var returnedAt string
		if val.ReturnedAt == nil {
			returnedAt = "not yet returned"
		} else {
			returnedAt = val.ReturnedAt.String()
		}
		overdueDto = append(overdueDto, models.OverdueTransactionDTO{
			ID:         val.ID.String(),
			BookID:     val.Book.ID.String(),
			BookName:   val.Book.Title,
			IssuedAt:   val.IssuedAt.String(),
			IssuedTill: val.IssuedTill.String(),
			ReturnedAt: returnedAt,
		})
	}

	return overdueDto, nil
}

func NewTransactionService(bookRepo bookrepo.BookStorage, transactionRepo transactionrepo.TransactionStorage) *TransactionService {
	return &TransactionService{
		bookRepo:        bookRepo,
		transactionRepo: transactionRepo,
	}
}

// IssueBook returns transaction id with error
func (service *TransactionService) IssueBook(ctx context.Context, bookId, issueFor string) (string, error) {
	userCtx, ok := ctx.Value("user").(models.UserJwt)
	if !ok {
		return "", errors.New("invalid user")
	}

	if userCtx.Role != roles.Customer {
		return "", errors.New("staff cant issue book")
	}

	if bookId == "" {
		return "", errors.New("invalid book id")
	}

	if issueFor == "" {
		issueFor = "1 day"
	}

	return service.transactionRepo.IssueBook(bookId, userCtx.ID, issueFor)
}

func (service *TransactionService) ReturnBook(ctx context.Context, bookId string) error {
	userCtx, ok := ctx.Value("user").(models.UserJwt)
	if !ok {
		return errors.New("invalid user")
	}

	if userCtx.Role != roles.Customer {
		return errors.New("staff cant return book")
	}

	if bookId == "" {
		return errors.New("invalid book id")
	}

	return service.transactionRepo.ReturnBook(bookId, userCtx.ID)
}

func (service *TransactionService) GetTransactions(ctx context.Context, dto models.GetTransactionRequestDTO) ([]models.TransactionDTO, error) {
	userCtx, ok := ctx.Value("user").(models.UserJwt)
	if !ok {
		return nil, errors.New("invalid user")
	}

	dto.UserId = userCtx.ID

	if dto.StartTime == "" {
		dto.StartTime = time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
	}

	if dto.EndTime == "" {
		dto.EndTime = time.Now().Format(time.RFC3339)
	}

	transactions, err := service.transactionRepo.GetAllTransactions(dto)
	if err != nil {
		return nil, err
	}

	var txDto []models.TransactionDTO
	for _, val := range transactions {

		var returnedAt string
		if val.ReturnedAt == nil {
			returnedAt = "not returned yet"
		} else {
			returnedAt = val.ReturnedAt.String()
		}

		txDto = append(txDto, models.TransactionDTO{
			ID:         val.ID.String(),
			BookID:     val.Book.ID.String(),
			BookName:   val.Book.Title,
			UserEmail:  val.User.Email,
			IssuedAt:   val.IssuedAt.String(),
			IssuedTill: val.IssuedTill.String(),
			ReturnedAt: returnedAt,
		})
	}

	return txDto, nil
}
