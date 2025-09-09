package transactionrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (repo *TransactionRepository) IssueBook(bookId, userId, issueFor string) (string, error) {
	var id string
	fmt.Println(issueFor)
	err := repo.db.QueryRow(`
		insert into transactions (book_id,user_id,issued_till)
		select $1, $2, now() + cast($3 as interval)
		where not exists(
		    select 1 from transactions where book_id = $1 and returned_at is null
		)
		returning id
`, bookId, userId, issueFor).Scan(&id)
	if err != nil {
		log.Println(err)
		return "", errors.New("book not available")
	}
	return id, err
}

func (repo *TransactionRepository) ReturnBook(bookId, userId string) error {
	res, err := repo.db.Exec(`update transactions set returned_at = $1 where book_id = $2 and user_id = $3 and returned_at is null`, time.Now(), bookId, userId)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("book already present nothing to return")
	}
	return nil
}

func (repo *TransactionRepository) GetAllTransactions(dto models.GetTransactionRequestDTO) ([]models.Transaction, error) {
	var transactions []models.Transaction

	rows, err := repo.db.Query(`
		select t.id, t.issued_at, t.returned_at, t.issued_till, b.id, b.title, u.email  from transactions as t
		left join users as u on t.user_id = u.id
		left join books as b on t.book_id = b.id
		where t.issued_at > $1 and t.issued_at < $2 
		and ($3='' or (returned_at is null) <> cast($3 as boolean))
		and ($4='' or b.title ilike '%'||$4||'%')
		and ($6='' or t.id = cast($6 as uuid))
		and u.id = $5
`, dto.StartTime, dto.EndTime, dto.Returned, dto.BookName, dto.UserId, dto.TransactionId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tx models.Transaction
		var returnedAt sql.Null[time.Time]
		err = rows.Scan(&tx.ID, &tx.IssuedAt, &returnedAt, &tx.IssuedTill, &tx.Book.ID, &tx.Book.Title, &tx.User.Email)
		if err != nil {
			return nil, err
		}

		if returnedAt.Valid {
			tx.ReturnedAt = &returnedAt.V
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (repo *TransactionRepository) GetOverDueTransactions(userId string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	rows, err := repo.db.Query(`
	select t.id, b.id, b.title, t.issued_at, t.issued_till, t.returned_at from transactions as t
	left join books as b
	on t.book_id = b.id
	where ($1='' or t.user_id = cast($1 as uuid) )
	and t.issued_till < now()
	and (t.returned_at is null or t.returned_at>t.issued_till)
`, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tx models.Transaction
		var returnedAt sql.Null[time.Time]
		err = rows.Scan(&tx.ID, &tx.Book.ID, &tx.Book.Title, &tx.IssuedAt, &tx.IssuedTill, &returnedAt)
		if err != nil {
			return nil, err
		}

		if returnedAt.Valid {
			tx.ReturnedAt = &returnedAt.V
		} else {
			tx.ReturnedAt = nil
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
