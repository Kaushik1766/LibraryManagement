package weberrors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendError(t *testing.T) {
	type args struct {
		err  error
		code int
		w    http.ResponseWriter
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "send error with 400 bad request",
			args: args{
				err:  errors.New("invalid request data"),
				code: http.StatusBadRequest,
				w:    httptest.NewRecorder(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid request data"}`,
		},
		{
			name: "send error with 401 unauthorized",
			args: args{
				err:  errors.New("authentication required"),
				code: http.StatusUnauthorized,
				w:    httptest.NewRecorder(),
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"authentication required"}`,
		},
		{
			name: "send error with 403 forbidden",
			args: args{
				err:  errors.New("access denied"),
				code: http.StatusForbidden,
				w:    httptest.NewRecorder(),
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"message":"access denied"}`,
		},
		{
			name: "send error with 404 not found",
			args: args{
				err:  errors.New("resource not found"),
				code: http.StatusNotFound,
				w:    httptest.NewRecorder(),
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"resource not found"}`,
		},
		{
			name: "send error with 500 internal server error",
			args: args{
				err:  errors.New("database connection failed"),
				code: http.StatusInternalServerError,
				w:    httptest.NewRecorder(),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database connection failed"}`,
		},
		{
			name: "send error with custom error message",
			args: args{
				err:  errors.New("book with ID 550e8400-e29b-41d4-a716-446655440000 not found"),
				code: http.StatusNotFound,
				w:    httptest.NewRecorder(),
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"book with ID 550e8400-e29b-41d4-a716-446655440000 not found"}`,
		},
		{
			name: "send error with empty error message",
			args: args{
				err:  errors.New(""),
				code: http.StatusBadRequest,
				w:    httptest.NewRecorder(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendError(tt.args.err, tt.args.code, tt.args.w)

			if recorder, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if recorder.Code != tt.expectedStatus {
					t.Errorf("SendError() status code = %v, want %v", recorder.Code, tt.expectedStatus)
				}

				body := strings.TrimSpace(recorder.Body.String())
				if body != tt.expectedBody {
					t.Errorf("SendError() body = %v, want %v", body, tt.expectedBody)
				}

				var webErr WebError
				if err := json.Unmarshal([]byte(body), &webErr); err != nil {
					t.Errorf("SendError() response is not valid JSON: %v", err)
				}

				if webErr.Message != tt.args.err.Error() {
					t.Errorf("SendError() message = %v, want %v", webErr.Message, tt.args.err.Error())
				}
			}
		})
	}
}

func TestWebError_JSONStructure(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple message",
			message:  "error occurred",
			expected: `{"message":"error occurred"}`,
		},
		{
			name:     "message with special characters",
			message:  "error: invalid input @ test.com",
			expected: `{"message":"error: invalid input @ test.com"}`,
		},
		{
			name:     "message with quotes",
			message:  `error: "field" is required`,
			expected: `{"message":"error: \"field\" is required"}`,
		},
		{
			name:     "empty message",
			message:  "",
			expected: `{"message":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webErr := WebError{Message: tt.message}
			data, err := json.Marshal(webErr)
			if err != nil {
				t.Errorf("Failed to marshal WebError: %v", err)
				return
			}

			result := string(data)
			if result != tt.expected {
				t.Errorf("WebError JSON = %v, want %v", result, tt.expected)
			}
		})
	}
}
