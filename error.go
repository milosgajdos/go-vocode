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

// TODO: APIError and APIAuthError are a single error but it's a Union

// APIParamError is returned when API params are invalid.
type APIParamError struct {
	Detail []struct {
		Loc  []string `json:"loc"`
		Msg  string   `json:"msg"`
		Type string   `json:"type"`
	} `json:"detail"`
}

// Error implements error interface.
func (e APIParamError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}

// APIGenError is returned when a generic API error is returned.
type APIGenError struct {
	Detail string `json:"detail"`
}

// Error implements error interface.
func (e APIGenError) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "unknown error"
	}
	return string(b)
}
