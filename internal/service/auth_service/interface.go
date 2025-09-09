package authservice

import "github.com/Kaushik1766/LibraryManagement/internal/models"

//go:generate mockgen -source=interface.go -destination=../../../mocks/mock_auth_manager.go -package=mocks
type AuthManager interface {
	Login(loginReq models.LoginDTO) (string, error)
	Signup(signupReq models.SignupDTO) error
}
