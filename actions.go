package vocode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

type TriggerType string

const (
	FnCallTriggerType TriggerType = "action_trigger_function_call"
	PhraseTriggerType TriggerType = "action_trigger_phrase_based"
)

type PhraseCondition string

const (
	PhraseTypeContains PhraseCondition = "phrase_condition_type_contains"
)

type Phrase struct {
	Phrase     string            `json:"phrase"`
	Conditions []PhraseCondition `json:"conditions"`
}

type PhraseTriggerConfig struct {
	PhraseTriggers []Phrase `json:"phrase_triggers"`
}

type PhraseTrigger struct {
	Type   TriggerType          `json:"type"`
	Config *PhraseTriggerConfig `json:"config"`
}

type FnCallTrigger struct {
	Type   TriggerType `json:"type"`
	Config map[any]any `json:"config"`
}

type ActionType string

const (
	TransferCall    ActionType = "action_transfer_call"
	EndConversation ActionType = "action_end_conversation"
	DTMF            ActionType = "action_dtmf"
	AddToConference ActionType = "action_add_to_conference"
	SetHold         ActionType = "action_set_hold"
	External        ActionType = "action_external"
)

type TransferCallActionConfig struct {
	PhoneNr string `json:"phone_number"`
}

type TransferCallAction struct {
	ActionBase
	Config *TransferCallActionConfig `json:"config"`
}

type EndConversationAction struct {
	ActionBase
	Config map[any]any `json:"config"`
}

type DTMFAction struct {
	ActionBase
	Config map[any]any `json:"config"`
}

type AddToConfConfig struct {
	PhoneNr            string `json:"phone_number"`
	PlacePrimaryOnHold bool   `json:"place_primary_on_hold"`
}

type AddToConfAction struct {
	ActionBase
	Config *AddToConfConfig `json:"config"`
}

type SetHoldAction struct {
	ActionBase
	Config map[any]any `json:"config"`
}

type ProcessingModeType string

const (
	MutedProcessingType ProcessingModeType = "muted"
)

type ExternalActionConfig struct {
	ProcessingMode ProcessingModeType `json:"processing_mode"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	URL            string             `json:"url"`
	InputSchema    map[any]any        `json:"input_schema"`
	SpeakOnSend    bool               `json:"speak_on_send"`
	SpeakonRecv    bool               `json:"speakon_recv"`
}

type ExternalAction struct {
	ActionBase
	Config *ExternalActionConfig `json:"config"`
}

type ActionBase struct {
	ID      string      `json:"id"`
	UserID  string      `json:"user_id"`
	Type    ActionType  `json:"type"`
	Trigger interface{} `json:"action_trigger"`
}

type Action struct {
	ActionBase
	Config interface{} `json:"config"`
}

type Actions struct {
	Items []Actions `json:"items"`
	*Paging
}

type CreateReq struct {
	Type    ActionType  `json:"type"`
	Config  interface{} `json:"config"`
	Trigger interface{} `json:"action_trigger"`
}

type UpdateReq struct {
	Type    ActionType  `json:"type"`
	Config  interface{} `json:"config"`
	Trigger interface{} `json:"action_trigger"`
}

func (a *Action) UnmarshalJSON(data []byte) error {
	type Alias Action
	aux := &struct {
		*Alias
		RawConfig      json.RawMessage `json:"config"`
		RawTrigger     json.RawMessage `json:"action_trigger"`
		RawTriggerType TriggerType     `json:"action_trigger_type"`
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.RawTriggerType {
	case FnCallTriggerType:
		var trigger FnCallTrigger
		if err := json.Unmarshal(aux.RawTrigger, &trigger); err != nil {
			return err
		}
		a.Trigger = &trigger
	case PhraseTriggerType:
		var trigger PhraseTrigger
		if err := json.Unmarshal(aux.RawTrigger, &trigger); err != nil {
			return err
		}
		a.Trigger = &trigger
	default:
		return fmt.Errorf("unknown trigger type: %s", aux.RawTriggerType)
	}

	switch a.Type {
	case TransferCall:
		var config TransferCallActionConfig
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = &config
	case EndConversation:
		var config map[string]interface{}
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = config
	case DTMF:
		var config map[string]interface{}
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = config
	case AddToConference:
		var config AddToConfConfig
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = &config
	case SetHold:
		var config map[string]interface{}
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = config
	case External:
		var config ExternalActionConfig
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = &config
	default:
		return fmt.Errorf("unknown action type: %s", a.Type)
	}

	return nil
}

func (c *Client) ListActions(ctx context.Context, paging *PageParams) (*Actions, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/actions/list")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}
	if paging != nil {
		request.WithPageParams(paging.Encode())
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		actions := new(Actions)
		if err := json.NewDecoder(resp.Body).Decode(actions); err != nil {
			return nil, err
		}
		return actions, nil
	case http.StatusForbidden:
		var apiErr APIAuthError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) GetAction(ctx context.Context, actionID string) (*Action, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/actions")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("id", actionID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Action)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden:
		var apiErr APIAuthError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) CreateAction(ctx context.Context, createReq *CreateReq) (*Action, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/actions/create")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(createReq); err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Action)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden:
		var apiErr APIAuthError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) UpdateAction(ctx context.Context, actionID string, updateReq *UpdateReq) (*Action, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/actions/update")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(updateReq); err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("id", actionID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Action)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden:
		var apiErr APIAuthError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}
