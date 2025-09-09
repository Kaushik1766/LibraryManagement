package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Kaushik1766/LibraryManagement/internal/config"
	"github.com/Kaushik1766/LibraryManagement/internal/models"
	weberrors "github.com/Kaushik1766/LibraryManagement/internal/web_errors"
	"github.com/golang-jwt/jwt/v5"
)

func TestParseToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    context.Context
		wantErr bool
	}{
		{
			name:    "invalid token",
			args:    args{token: "adfafa"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong signature",
			args: args{
				token: func() string {
					claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserJwt{
						Email: "kaushik@a.com",
						Role:  0,
					})

					signedToken, _ := claims.SignedString([]byte("ddafs"))
					return signedToken
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "expired token",
			args: args{
				token: func() string {
					claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserJwt{
						Email: "kaushik@a.com",
						Role:  0,
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * -10)),
						},
					})

					signedToken, _ := claims.SignedString([]byte(config.JWTSecret))
					return signedToken
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid token",
			args: args{
				token: func() string {
					claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserJwt{
						Email: "kaushik@a.com",
						Role:  0,
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Date(10000, 1, 1, 1, 1, 1, 1, time.Local)),
						},
					})

					signedToken, _ := claims.SignedString([]byte(config.JWTSecret))
					return signedToken
				}(),
			},
			want: context.WithValue(context.Background(), "user", models.UserJwt{
				Email: "kaushik@a.com",
				Role:  0,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Date(10000, 1, 1, 1, 1, 1, 1, time.Local)),
				},
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	createValidToken := func() string {
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserJwt{
			Email: "test@example.com",
			Role:  0,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})
		token, _ := claims.SignedString([]byte(config.JWTSecret))
		return token
	}

	createExpiredToken := func() string {
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserJwt{
			Email: "test@example.com",
			Role:  0,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			},
		})
		token, _ := claims.SignedString([]byte(config.JWTSecret))
		return token
	}

	createInvalidToken := func() string {
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserJwt{
			Email: "test@example.com",
			Role:  0,
		})
		token, _ := claims.SignedString([]byte("wrong-secret"))
		return token
	}

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		nextCalled     bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "missing authorization header",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			nextCalled:     false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name: "invalid authorization header format - no Bearer prefix",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", createValidToken())
				return req
			},
			nextCalled:     false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name: "invalid authorization header format - malformed Bearer",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer")
				return req
			},
			nextCalled:     false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name: "invalid JWT token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer invalid-token")
				return req
			},
			nextCalled:     false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "token is malformed: token contains an invalid number of segments",
		},
		{
			name: "JWT token with wrong signature",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+createInvalidToken())
				return req
			},
			nextCalled:     false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "token signature is invalid: signature is invalid",
		},
		{
			name: "expired JWT token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+createExpiredToken())
				return req
			},
			nextCalled:     false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "token has invalid claims: token is expired",
		},
		{
			name: "valid JWT token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+createValidToken())
				return req
			},
			nextCalled:     true,
			expectedStatus: http.StatusOK,
			expectedBody:   "next handler called",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := tt.setupRequest()

			nextCalled := false
			var capturedCtx context.Context

			nextHandler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				capturedCtx = ctx
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("next handler called"))
			}

			middleware := AuthMiddleware(nextHandler)
			middleware(w, r)

			if nextCalled != tt.nextCalled {
				t.Errorf("AuthMiddleware() next handler called = %v, want %v", nextCalled, tt.nextCalled)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("AuthMiddleware() status code = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedBody != "" && !tt.nextCalled {
				var webErr weberrors.WebError
				if err := json.Unmarshal(w.Body.Bytes(), &webErr); err != nil {
					t.Errorf("AuthMiddleware() failed to unmarshal error response: %v", err)
				}
				if webErr.Message != tt.expectedBody {
					t.Errorf("AuthMiddleware() error message = %v, want %v", webErr.Message, tt.expectedBody)
				}
			} else if tt.nextCalled {
				body := w.Body.String()
				if body != tt.expectedBody {
					t.Errorf("AuthMiddleware() response body = %v, want %v", body, tt.expectedBody)
				}
			}

			if tt.nextCalled && tt.name == "valid JWT token" {
				userValue := capturedCtx.Value("user")
				if userValue == nil {
					t.Error("AuthMiddleware() context should contain user info")
				} else {
					userJwt, ok := userValue.(models.UserJwt)
					if !ok {
						t.Error("AuthMiddleware() user value should be of type UserJwt")
					} else if userJwt.Email != "test@example.com" {
						t.Errorf("AuthMiddleware() user email = %v, want test@example.com", userJwt.Email)
					}
				}
			}
		})
	}
}
