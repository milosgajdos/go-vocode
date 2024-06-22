package vocode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

// ActionTrigger must be implemented by
// all action triggers.
type ActionTrigger interface {
	TriggerType() TriggerType
}

type TriggerType string

const (
	FnCallTriggerType TriggerType = "action_trigger_function_call"
	PhraseTriggerType TriggerType = "action_trigger_phrase_based"
)

type PhraseCondition string

const (
	PhraseCondTypeContains PhraseCondition = "phrase_condition_type_contains"
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

func (p *PhraseTrigger) TriggerType() TriggerType {
	return PhraseTriggerType
}

type FnCallTrigger struct {
	Type   TriggerType    `json:"type"`
	Config map[string]any `json:"config"`
}

func (f *FnCallTrigger) TriggerType() TriggerType {
	return FnCallTriggerType
}

type ActionType string

const (
	ActionDTMF            ActionType = "action_dtmf"
	ActionSetHold         ActionType = "action_set_hold"
	ActionExternal        ActionType = "action_external"
	ActionTransferCall    ActionType = "action_transfer_call"
	ActionEndConversation ActionType = "action_end_conversation"
	ActionAddToConference ActionType = "action_add_to_conference"
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
	Config map[string]any `json:"config"`
}

type DTMFAction struct {
	ActionBase
	Config map[string]any `json:"config"`
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
	Config map[string]any `json:"config"`
}

type ProcessingMode string

const (
	MutedProcessing ProcessingMode = "muted"
)

type ExternalActionConfig struct {
	ProcessingMode ProcessingMode `json:"processing_mode"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	URL            string         `json:"url"`
	InputSchema    map[string]any `json:"input_schema"`
	SpeakOnSend    bool           `json:"speak_on_send"`
	SpeakonRecv    bool           `json:"speak_on_receive"`
}

type ExternalAction struct {
	ActionBase
	Config *ExternalActionConfig `json:"config"`
}

type ActionBase struct {
	ID      string        `json:"id"`
	UserID  string        `json:"user_id"`
	Type    ActionType    `json:"type"`
	Trigger ActionTrigger `json:"action_trigger"`
}

type Action struct {
	ActionBase
	Config interface{} `json:"config,omitempty"`
}

type Actions struct {
	Items []Action `json:"items"`
	*Paging
}

type ActionReqBase struct {
	Type    ActionType    `json:"type"`
	Trigger ActionTrigger `json:"action_trigger"`
	Config  interface{}   `json:"config"`
}

type CreateActionReq struct {
	ActionReqBase
}

type UpdateActionReq struct {
	ActionReqBase
}

func (a *Action) UnmarshalJSON(data []byte) error {
	type Alias Action
	aux := &struct {
		*Alias
		RawConfig  json.RawMessage `json:"config"`
		RawTrigger json.RawMessage `json:"action_trigger"`
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var triggerObj struct {
		Type   TriggerType     `json:"type"`
		Config json.RawMessage `json:"config"`
	}

	if err := json.Unmarshal(aux.RawTrigger, &triggerObj); err != nil {
		return err
	}

	switch triggerObj.Type {
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
		return fmt.Errorf("unknown trigger type: %s", triggerObj.Type)
	}

	switch a.Type {
	case ActionTransferCall:
		var config TransferCallActionConfig
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = &config
	case ActionEndConversation:
		var config map[string]interface{}
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = config
	case ActionDTMF:
		var config map[string]interface{}
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = config
	case ActionAddToConference:
		var config AddToConfConfig
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = &config
	case ActionSetHold:
		var config map[string]interface{}
		if err := json.Unmarshal(aux.RawConfig, &config); err != nil {
			return err
		}
		a.Config = config
	case ActionExternal:
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

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	actions := new(Actions)
	if err := json.NewDecoder(resp.Body).Decode(actions); err != nil {
		return nil, err
	}
	return actions, nil
}

func (c *Client) GetAction(ctx context.Context, id string) (*Action, error) {
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
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	action := new(Action)
	if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
		return nil, err
	}
	return action, nil
}

func (c *Client) CreateAction(ctx context.Context, createReq *CreateActionReq) (*Action, error) {
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

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	action := new(Action)
	if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
		return nil, err
	}
	return action, nil
}

func (c *Client) UpdateAction(ctx context.Context, id string, updateReq *UpdateActionReq) (*Action, error) {
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
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	action := new(Action)
	if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
		return nil, err
	}
	return action, nil
}
