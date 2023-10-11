package api

import (
	"encoding/json"
)

type Payload interface {
	Validate() *ValidationError
}

type ValidationDetails map[string]string

type ValidationError struct {
	Details ValidationDetails `json:"details"`
}

func (e ValidationError) Error() string {
	r, _ := json.Marshal(e.Details)
	return string(r)
}

type ErrorResponse struct {
	Error  string            `json:"error"`
	Detail map[string]string `json:"detail,omitempty"`
}

type APIError struct {
	StatusCode int               `json:"-"`
	Status     string            `json:"error"`
	Detail     map[string]string `json:"detail,omitempty"`
}

func NewAPIError(statusCode int, status string, detail map[string]string) APIError {
	return APIError{statusCode, status, detail}
}

func (e APIError) Error() string {
	return e.Status
}
