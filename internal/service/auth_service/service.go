package authservice

import (
	"errors"
	"net/mail"
	"time"

	"github.com/Kaushik1766/LibraryManagement/internal/config"
	"github.com/Kaushik1766/LibraryManagement/internal/models"
	userrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/user_repo"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo userrepo.UserStorage
}

func NewAuthService(userRepo userrepo.UserStorage) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (service *AuthService) Login(loginReq models.LoginDTO) (string, error) {
	if loginReq.Email == "" || loginReq.Password == "" {
		return "", errors.New("email or password cant be empty")
	}
	_, err := mail.ParseAddress(loginReq.Email)
	if err != nil {
		return "", errors.New("invalid email address")
	}

	user, err := service.userRepo.GetUserByEmail(loginReq.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserJwt{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: user.Email,
		Role:  user.Role,
	})

	return token.SignedString([]byte(config.JWTSecret))
}

func (service *AuthService) Signup(signupReq models.SignupDTO) error {
	if signupReq.Name == "" || signupReq.Password == "" || signupReq.Email == "" {
		return errors.New("name, email or password cant be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), 12)
	if err != nil {
		return errors.New("password too long")
	}

	return service.userRepo.AddUser(signupReq.Name, signupReq.Email, string(hashedPassword))
}
