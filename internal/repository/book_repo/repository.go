package bookrepo

import (
	"database/sql"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (repo *BookRepository) AddBook(title, author string, copies int) error {
	_, err := repo.db.Exec(
		`insert into books(title, author) 
		select $1, $2
		from generate_series(1,$3)`,
		title, author, copies)
	return err
}

func (repo *BookRepository) GetAllBooks(title, author string) ([]models.Book, error) {
	var books []models.Book
	rows, err := repo.db.Query(`
	select b.id, b.title, b.author, u.email from books as b left join transactions as t
	on t.id = (
	    select t1.id from transactions as t1
	                where t1.book_id = b.id
					order by t1.issued_at desc 
					limit 1
	)
	left join users as u
	on t.user_id = u.id and t.returned_at is null
	where ($1='' or b.title ilike '%'||$1||'%')
	and ($2='' or b.author ilike '%'||$2||'%')
`, title, author)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b models.Book
		b.IssuedTo = &models.User{}
		var email sql.NullString
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &email)
		if err != nil {
			return nil, err
		}

		if email.Valid {
			b.IssuedTo.Email = email.String
		} else {
			b.IssuedTo = nil
		}
		books = append(books, b)
	}
	return books, nil
}
