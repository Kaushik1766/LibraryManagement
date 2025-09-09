package authservice

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Kaushik1766/LibraryManagement/internal/config"
	"github.com/Kaushik1766/LibraryManagement/internal/models"
	userrepo "github.com/Kaushik1766/LibraryManagement/internal/repository/user_repo"
	"github.com/Kaushik1766/LibraryManagement/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)

	type fields struct {
		userRepo userrepo.UserStorage
	}
	type args struct {
		loginReq models.LoginDTO
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		checkOutput func(string) bool
		wantErr     bool
		mockSetup   func()
	}{
		{
			name:   "valid login",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				loginReq: models.LoginDTO{
					Email:    "kaushik@a.com",
					Password: "123",
				},
			},
			checkOutput: func(token string) bool {
				var userJwt models.UserJwt
				jwtToken, err := jwt.ParseWithClaims(token, &userJwt, func(token *jwt.Token) (any, error) {
					return []byte(config.JWTSecret), nil
				})

				if err != nil {
					return false
				}

				if !jwtToken.Valid {
					return false
				}

				if userJwt.Email != "kaushik@a.com" {
					return false
				}

				return true
			},
			wantErr: false,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByEmail("kaushik@a.com").Return(models.User{
					ID:   uuid.New(),
					Name: "kaushik",
					Password: func() string {
						hash, _ := bcrypt.GenerateFromPassword([]byte("123"), 12)
						return string(hash)
					}(),
					Email: "kaushik@a.com",
					Role:  0,
				}, nil)
			},
		},
		{
			name:   "invalid login",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				loginReq: models.LoginDTO{
					Email:    "kaushik@a.com",
					Password: "123",
				},
			},
			checkOutput: func(token string) bool {
				return true
			},
			wantErr: true,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByEmail("kaushik@a.com").Return(models.User{}, errors.New("db error"))
			},
		},
		{
			name:   "wrong password",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				loginReq: models.LoginDTO{
					Email:    "kaushik@a.com",
					Password: "125",
				},
			},
			checkOutput: func(token string) bool {
				return true
			},
			wantErr: true,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserByEmail("kaushik@a.com").Return(models.User{
					ID:   uuid.New(),
					Name: "kaushik",
					Password: func() string {
						hash, _ := bcrypt.GenerateFromPassword([]byte("123"), 12)
						return string(hash)
					}(),
					Email: "kaushik@a.com",
					Role:  0,
				}, nil)
			},
		},
		{
			name:   "invalid email",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				loginReq: models.LoginDTO{
					Email:    "kaushika.com",
					Password: "125",
				},
			},
			checkOutput: func(token string) bool {
				return true
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name:   "incomplete fields",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				loginReq: models.LoginDTO{
					Email:    "",
					Password: "125",
				},
			},
			checkOutput: func(token string) bool {
				return true
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AuthService{
				userRepo: tt.fields.userRepo,
			}
			tt.mockSetup()
			got, err := service.Login(tt.args.loginReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.checkOutput(got) {
				t.Errorf("invalid jwt")
			}
		})
	}
}

func TestAuthService_Signup(t *testing.T) {

	ctrl := gomock.NewController(t)

	mockUserRepo := mocks.NewMockUserStorage(ctrl)

	type fields struct {
		userRepo userrepo.UserStorage
	}
	type args struct {
		signupReq models.SignupDTO
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name:   "valid signup",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				signupReq: models.SignupDTO{
					Name:     "kaushik",
					Email:    "kaushik@a.com",
					Password: "123",
				},
			},
			wantErr: false,
			mockSetup: func() {
				mockUserRepo.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:   "invalid signup",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				signupReq: models.SignupDTO{
					Name:     "",
					Email:    "kaushik@a.com",
					Password: "123",
				},
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
		{
			name:   "password too long",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				signupReq: models.SignupDTO{
					Name:  "kaushik",
					Email: "kaushik@a.com",
					Password: func() string {
						var s string
						for range 100 {
							s += "a"
						}
						return s
					}(),
				},
			},
			wantErr: true,
			mockSetup: func() {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AuthService{
				userRepo: tt.fields.userRepo,
			}
			tt.mockSetup()
			if err := service.Signup(tt.args.signupReq); (err != nil) != tt.wantErr {
				t.Errorf("Signup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	type args struct {
		userRepo userrepo.UserStorage
	}
	tests := []struct {
		name string
		args args
		want *AuthService
	}{
		{
			name: "valid",
			args: args{userRepo: mockUserRepo},
			want: &AuthService{userRepo: mockUserRepo},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthService(tt.args.userRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}
