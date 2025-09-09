package models

import (
	"github.com/Kaushik1766/LibraryManagement/internal/models/enums/roles"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Password string
	Email    string
	Role     roles.UserRoles
}

type UserJwt struct {
	jwt.RegisteredClaims
	Email string
	Role  roles.UserRoles
}

type SignupDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
