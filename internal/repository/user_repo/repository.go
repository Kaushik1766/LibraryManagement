package userrepo

import (
	"database/sql"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u UserRepository) AddUser(name, email, password string) error {
	_, err := u.db.Exec(`insert into users(name, email, password) values($1,$2,$3)`, name, email, password)
	return err
}

func (u UserRepository) GetUserByEmail(email string) (models.User, error) {
	row := u.db.QueryRow(`select * from users where email = $1`, email)

	var user models.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	return user, err
}
