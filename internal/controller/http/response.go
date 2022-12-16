package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var RequiredFieldMsg = "обязательное поле"

// httpResponse is base HTTPResponse struct view.
type httpResponse struct {
	Status       string         `json:"status,omitempty"`
	Code         int            `json:"status_code,omitempty"`
	Message      string         `json:"message,omitempty"`
	Payload      PayloadList    `json:"payload,omitempty"`
	Validation   ValidationList `json:"validation,omitempty"`
	Debug        string         `json:"debug,omitempty"`
	JWT          string         `json:"jwt_token,omitempty"`
	RefreshToken string         `json:"refresh_token,omitempty"`
}

type (
	PayloadList   []map[string]interface{}
	payloadOption func(*httpResponse)

	ValidationList   []map[string]interface{}
	validationOption func(*httpResponse)
)

func withPayload(key string, value interface{}) payloadOption {
	return func(hr *httpResponse) {
		pl := make(map[string]interface{}, 1)
		pl[key] = value
		hr.Payload = append(hr.Payload, pl)
	}
}

func withValidation(key string) validationOption {
	return func(hr *httpResponse) {
		pl := make(map[string]interface{}, 1)
		pl[key] = key + " " + RequiredFieldMsg
		hr.Validation = append(hr.Validation, pl)
	}
}

func successResponse(c echo.Context, code int, msg string, opts ...payloadOption) error {
	response := &httpResponse{
		Code:    code,
		Message: msg,
		Status:  http.StatusText(code),
	}

	for _, opt := range opts {
		opt(response)
	}

	return c.JSON(code, response)
}

func errorResponse(c echo.Context, code int, msg string, err error, fields []string) error {
	response := &httpResponse{
		Code:    code,
		Message: msg,
		Status:  http.StatusText(code),
		Debug:   err.Error(),
	}

	for _, field := range fields {
		f := withValidation(field)
		f(response)
	}

	return c.JSON(code, response)
}
