package transactionhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	transactionservice "github.com/Kaushik1766/LibraryManagement/internal/service/transaction_service"
	weberrors "github.com/Kaushik1766/LibraryManagement/internal/web_errors"
)

type TransactionHandler struct {
	transactionService transactionservice.TransactionManager
}

func NewTransactionHandler(transactionService transactionservice.TransactionManager) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (handler *TransactionHandler) IssueBook(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req struct {
		BookId    string `json:"book_id"`
		IssueTill string `json:"issue_for"`
	}

	data, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(data, &req)
	if err != nil {
		weberrors.SendError(err, http.StatusBadRequest, w)
		return
	}

	transactionId, err := handler.transactionService.IssueBook(ctx, req.BookId, req.IssueTill)
	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"transaction_id":"%s"}`, transactionId)))
}

func (handler *TransactionHandler) ReturnBook(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req struct {
		BookId string `json:"book_id"`
	}

	data, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(data, &req)
	if err != nil {
		weberrors.SendError(err, http.StatusBadRequest, w)
		return
	}

	err = handler.transactionService.ReturnBook(ctx, req.BookId)
	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *TransactionHandler) GetAllTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	transactions, err := handler.transactionService.GetTransactions(ctx, models.GetTransactionRequestDTO{
		UserId:    "",
		StartTime: query.Get("startTime"),
		EndTime:   query.Get("endTime"),
		Returned:  query.Get("returned"),
		BookName:  query.Get("title"),
	})

	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

func (handler *TransactionHandler) GetOverdueTransactions(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	overdueTransactions, err := handler.transactionService.GetOverdueTransactions(ctx)
	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(overdueTransactions)
}

func (handler *TransactionHandler) GetTransactionById(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	transactionId := r.PathValue("transactionId")

	transactions, err := handler.transactionService.GetTransactions(ctx, models.GetTransactionRequestDTO{
		TransactionId: transactionId,
	})

	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}
