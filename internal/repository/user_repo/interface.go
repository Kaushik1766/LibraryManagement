package userrepo

import "github.com/Kaushik1766/LibraryManagement/internal/models"

//go:generate mockgen -source=interface.go -destination=../../../mocks/mock_user_storage.go -package=mocks
type UserStorage interface {
	AddUser(name, email, password string) error
	GetUserByEmail(email string) (models.User, error)
}
