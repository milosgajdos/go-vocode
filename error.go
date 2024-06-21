package vocode

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	ParamError     *APIParamError
	GenError       *APIGenError
	UnexpecedError json.RawMessage
}

func (e *APIError) Error() string {
	if e.ParamError != nil {
		return e.ParamError.Error()
	}
	if e.GenError != nil {
		return e.GenError.Error()
	}
	if len(e.UnexpecedError) > 0 {
		return string(e.UnexpecedError)
	}
	return "unknown error"
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	var paramError APIParamError
	err := json.Unmarshal(data, &paramError)
	if err == nil && len(paramError.Detail) > 0 {
		e.ParamError = &paramError
		return nil
	}

	var genError APIGenError
	err = json.Unmarshal(data, &genError)
	if err == nil && genError.Detail != "" {
		e.GenError = &genError
		return nil
	}

	e.UnexpecedError = data

	return fmt.Errorf("unexpected API error JSON: %s", string(data))
}

// APIParamError is returned when API params are invalid.
type APIParamError struct {
	Detail []struct {
		Type  string        `json:"type"`
		Loc   []interface{} `json:"loc"`
		Msg   string        `json:"msg"`
		Input string        `json:"input,omitempty"`
		Ctx   *struct {
			Error string `json:"error"`
		} `json:"ctx,omitempty"`
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
