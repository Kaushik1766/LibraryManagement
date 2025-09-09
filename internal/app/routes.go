package app

import (
	"net/http"

	"github.com/Kaushik1766/LibraryManagement/internal/middleware"
)

var routes map[string]func(w http.ResponseWriter, r *http.Request)

func (app *App) registerRoutes() {

	authMiddleware := middleware.AuthMiddleware

	routes = map[string]func(w http.ResponseWriter, r *http.Request){
		"POST /auth/signup":                 app.AuthHandler.Signup,
		"POST /auth/login":                  app.AuthHandler.Login,
		"POST /books":                       authMiddleware(app.BookHandler.AddBook),
		"GET /books":                        authMiddleware(app.BookHandler.GetAllBooks),
		"POST /transactions/issue":          authMiddleware(app.TransactionHandler.IssueBook),
		"POST /transactions/return":         authMiddleware(app.TransactionHandler.ReturnBook),
		"GET /transactions/overdue":         authMiddleware(app.TransactionHandler.GetOverdueTransactions),
		"GET /transactions":                 authMiddleware(app.TransactionHandler.GetAllTransactions),
		"GET /transactions/{transactionId}": authMiddleware(app.TransactionHandler.GetTransactionById),
	}

	for route, handler := range routes {
		app.mux.HandleFunc(route, handler)
	}
}
