package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Kaushik1766/LibraryManagement/internal/config"
	"github.com/Kaushik1766/LibraryManagement/internal/models"
	weberrors "github.com/Kaushik1766/LibraryManagement/internal/web_errors"
	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(token string) (context.Context, error) {

	var userJwt models.UserJwt

	parsedToken, err := jwt.ParseWithClaims(token, &userJwt, func(token *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	if userJwt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return context.WithValue(context.Background(), "user", userJwt), nil
}

func AuthMiddleware(next func(ctx context.Context, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(token) <= 7 || token[:7] != "Bearer " {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtToken := token[7:]

		ctx, err := ParseToken(jwtToken)
		if err != nil {
			weberrors.SendError(err, http.StatusUnauthorized, w)
			return
		}

		next(ctx, w, r)
	}
}
