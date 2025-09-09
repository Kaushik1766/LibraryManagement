package app

import (
	"database/sql"
	"fmt"
	"net/http"

	authhandler "github.com/Kaushik1766/LibraryManagement/internal/handlers/auth_handler"
	"github.com/Kaushik1766/LibraryManagement/internal/handlers/book_handler"
	transactionhandler "github.com/Kaushik1766/LibraryManagement/internal/handlers/transaction_handler"
	bookrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/book_repo"
	transactionrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/transaction_repo"
	userrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/user_repo"
	authservice "github.com/Kaushik1766/LibraryManagement/internal/service/auth_service"
	bookservice "github.com/Kaushik1766/LibraryManagement/internal/service/book_service"
	transactionservice "github.com/Kaushik1766/LibraryManagement/internal/service/transaction_service"
)

var (
	userRepo        userrepo.UserStorage               = nil
	bookRepo        bookrepo.BookStorage               = nil
	transactionRepo transactionrepo.TransactionStorage = nil

	authService        authservice.AuthManager               = nil
	bookService        bookservice.BookManager               = nil
	transactionService transactionservice.TransactionManager = nil
)

type App struct {
	mux *http.ServeMux
	db  *sql.DB

	AuthHandler        *authhandler.AuthHandler
	BookHandler        *bookhandler.BookHandler
	TransactionHandler *transactionhandler.TransactionHandler
}

func NewApp(db *sql.DB) *App {
	app := App{
		mux: http.NewServeMux(),
		db:  db,
	}

	userRepo = userrepo.NewUserRepository(db)
	bookRepo = bookrepo.NewBookRepository(db)
	transactionRepo = transactionrepo.NewTransactionRepository(db)

	authService = authservice.NewAuthService(userRepo)
	bookService = bookservice.NewBookService(bookRepo)
	transactionService = transactionservice.NewTransactionService(bookRepo, transactionRepo)

	app.AuthHandler = authhandler.NewAuthHandler(authService)
	app.BookHandler = bookhandler.NewBookHandler(bookService)
	app.TransactionHandler = transactionhandler.NewTransactionHandler(transactionService)

	app.registerRoutes()
	return &app
}

func (app *App) Run() {
	fmt.Println("server started at port 3000")
	http.ListenAndServe("localhost:3000", app.mux)
}
