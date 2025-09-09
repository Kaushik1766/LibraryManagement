package weberrors

import (
	"encoding/json"
	"net/http"
)

type WebError struct {
	Message string `json:"message"`
}

func SendError(err error, code int, w http.ResponseWriter) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(WebError{
		Message: err.Error(),
	})
}
