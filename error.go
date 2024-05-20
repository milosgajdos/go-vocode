package vocode

import (
	"encoding/json"
	"errors"
)

var (
	// ErrTooManyRequests is returned when the cient hits rate limit.
	ErrTooManyRequests = errors.New("too many requests")
	// ErrUnexpectedStatusCode is returned when an unexpected status is returned from the API.
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	// ErrUnprocessableEntity is returned when no valid data is available
	ErrUnprocessableEntity = errors.New("unprocessable entity")
)

// APIError is open AI API error.
type APIError struct {
	Detail []struct {
		Loc  []string `json:"loc"`
		Msg  string   `json:"msg"`
		Type string   `json:"type"`
	} `json:"detail"`
}

// Error implements error interface.
func (e APIError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}

// APIAuthError is returned when API auth fails.
type APIAuthError struct {
	Detail string `json:"detail"`
}

// Error implements error interface.
func (e APIAuthError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}
