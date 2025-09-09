package authhandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	authservice "github.com/Kaushik1766/LibraryManagement/internal/service/auth_service"
	"github.com/Kaushik1766/LibraryManagement/mocks"
	"go.uber.org/mock/gomock"
)

func anyToReader(data any) io.Reader {
	dataJsonBytes, _ := json.Marshal(data)
	return bytes.NewReader(dataJsonBytes)
}

func TestAuthHandler_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	authService := mocks.NewMockAuthManager(ctrl)

	type fields struct {
		authService authservice.AuthManager
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		checkOutput func(recorder *httptest.ResponseRecorder) error
		mockSetup   func()
	}{
		{
			name:   "valid login",
			fields: fields{authService: authService},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/login", anyToReader(models.LoginDTO{
					Email:    "kaushik@a.com",
					Password: "123",
				})),
			},
			checkOutput: func(recorder *httptest.ResponseRecorder) error {
				resp := recorder.Result()
				if resp.StatusCode != http.StatusOK {
					return errors.New("wrong status code")
				}
				data, _ := io.ReadAll(resp.Body)
				var token struct {
					Jwt string `json:"jwt"`
				}

				err := json.Unmarshal(data, &token)
				if err != nil {
					return errors.New("invalid response body: " + err.Error())
				}

				if token.Jwt == "validToken" {
					return nil
				} else {
					return errors.New("invalid token")
				}
			},
			mockSetup: func() {
				authService.EXPECT().Login(gomock.Any()).Return("validToken", nil)
			},
		},
		{
			name:   "invalid json",
			fields: fields{authService: authService},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/login", anyToReader("invalid json")),
			},
			checkOutput: func(recorder *httptest.ResponseRecorder) error {
				resp := recorder.Result()
				if resp.StatusCode != http.StatusBadRequest {
					return errors.New("wrong status code")
				}
				return nil
			},
			mockSetup: func() {
			},
		},
		{
			name:   "failed login",
			fields: fields{authService: authService},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/login", anyToReader(models.LoginDTO{
					Email:    "kaushik@a.com",
					Password: "123",
				})),
			},
			checkOutput: func(recorder *httptest.ResponseRecorder) error {
				resp := recorder.Result()
				if resp.StatusCode != http.StatusInternalServerError {
					return errors.New("wrong status code")
				}
				return nil
			},
			mockSetup: func() {
				authService.EXPECT().Login(gomock.Any()).Return("", errors.New("invalid credentials"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &AuthHandler{
				authService: tt.fields.authService,
			}
			tt.mockSetup()
			handler.Login(tt.args.w, tt.args.r)
			if err := tt.checkOutput(tt.args.w); err != nil {
				t.Errorf("invalid output %s", err.Error())
			}
		})
	}
}

func TestAuthHandler_Signup(t *testing.T) {

	ctrl := gomock.NewController(t)
	authService := mocks.NewMockAuthManager(ctrl)

	type fields struct {
		authService authservice.AuthManager
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		checkOutput func(recorder *httptest.ResponseRecorder) error
		mockSetup   func()
	}{
		{
			name:   "invalid json",
			fields: fields{authService: authService},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/signup", anyToReader("invalid json")),
			},
			checkOutput: func(recorder *httptest.ResponseRecorder) error {
				resp := recorder.Result()
				if resp.StatusCode != http.StatusBadRequest {
					return errors.New("invalid response code")
				}
				return nil
			},
			mockSetup: func() {
			},
		},
		{
			name:   "valid signup",
			fields: fields{authService: authService},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/signup", anyToReader(models.SignupDTO{
					Name:     "kaushik",
					Email:    "kaushik@a.com",
					Password: "123",
				})),
			},
			checkOutput: func(recorder *httptest.ResponseRecorder) error {
				resp := recorder.Result()
				if resp.StatusCode != http.StatusOK {
					return errors.New("invalid response code")
				}
				return nil
			},
			mockSetup: func() {
				authService.EXPECT().Signup(gomock.Any()).Return(nil)
			},
		},
		{
			name:   "invalid signup",
			fields: fields{authService: authService},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/signup", anyToReader(models.SignupDTO{
					Name:     "kaushik",
					Email:    "kaushik@a.com",
					Password: "123",
				})),
			},
			checkOutput: func(recorder *httptest.ResponseRecorder) error {
				resp := recorder.Result()
				if resp.StatusCode != http.StatusInternalServerError {
					return errors.New("invalid response code")
				}
				return nil
			},
			mockSetup: func() {
				authService.EXPECT().Signup(gomock.Any()).Return(errors.New("user exists"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &AuthHandler{
				authService: tt.fields.authService,
			}
			tt.mockSetup()
			handler.Signup(tt.args.w, tt.args.r)
			if err := tt.checkOutput(tt.args.w); err != nil {
				t.Errorf("invalid repsonse: %s", err)
			}
		})
	}
}

func TestNewAuthHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	authService := mocks.NewMockAuthManager(ctrl)
	type args struct {
		authService authservice.AuthManager
	}
	tests := []struct {
		name string
		args args
		want *AuthHandler
	}{
		{
			name: "valid",
			args: args{authService},
			want: &AuthHandler{authService},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthHandler(tt.args.authService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
