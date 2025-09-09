package authhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Kaushik1766/LibraryManagement/internal/models"
	authservice "github.com/Kaushik1766/LibraryManagement/internal/service/auth_service"
	weberrors "github.com/Kaushik1766/LibraryManagement/internal/web_errors"
)

type AuthHandler struct {
	authService authservice.AuthManager
}

func NewAuthHandler(authService authservice.AuthManager) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (handler *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.SignupDTO
	data, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(data, &req)
	if err != nil {
		weberrors.SendError(err, http.StatusBadRequest, w)
		return
	}
	err = handler.authService.Signup(req)
	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginDTO
	data, _ := io.ReadAll(r.Body)

	err := json.Unmarshal(data, &req)
	if err != nil {
		weberrors.SendError(err, http.StatusBadRequest, w)
		return
	}
	token, err := handler.authService.Login(req)
	if err != nil {
		weberrors.SendError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"jwt":"%s"}`, token)))
}
