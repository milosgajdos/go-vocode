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

type Language string

const (
	English    Language = "en"
	Spanish    Language = "es"
	German     Language = "de"
	Portuguese Language = "pt"
	French     Language = "fr"
	Hindi      Language = "hi"
	Dutch      Language = "nl"
	Italian    Language = "it"
	Japanese   Language = "jp"
	Korean     Language = "ko"
)

type InterruptSenseType string

const (
	LowInterruptSense  InterruptSenseType = "low"
	HighInterruptSense InterruptSenseType = "high"
)

type EndpointSenseType string

const (
	AutoEndpointSense      EndpointSenseType = "auto"
	RelaxedEndpointSense   EndpointSenseType = "relaxed"
	SensitiveEndpointSense EndpointSenseType = "sensitive"
)

type IVRNavModeType string

const (
	DefaultIVRMode IVRNavModeType = "default"
	OffIVRMode     IVRNavModeType = "off"
)

type AgentOpenAIAccount struct {
	AccountConnsBase
	OpenAIAccount *OpenAIAccount
}

type Agents struct {
	Items []Agent `json:"items"`
	*Paging
}

type Agent struct {
	ID                       string              `json:"id"`
	UserID                   string              `json:"user_id"`
	Name                     string              `json:"name"`
	Prompt                   *Prompt             `json:"prompt"`
	Language                 Language            `json:"language"`
	Actions                  []Action            `json:"actions"`
	Voice                    *Voice              `json:"voice"`
	InitMsg                  string              `json:"initial_msg"`
	Webhook                  *Webhook            `json:"webhook"`
	VectorDB                 *VectorDB           `json:"vector_database"`
	InterruptSense           InterruptSenseType  `json:"interrupt_sensitivity"`
	CtxEndpint               string              `json:"context_endpint"`
	NoiseSuppression         bool                `json:"noise_suppression"`
	EndpointSense            EndpointSenseType   `json:"endpointing_sensitivity"`
	IVRNavMode               IVRNavModeType      `json:"ivr_navigation_mode"`
	Speed                    float32             `json:"conversation_speed"`
	InitMsgDelay             float64             `json:"initial_message_delay"`
	OpenAIModelOverride      bool                `json:"openai_model_name_override"`
	AsktIfHumanPresentOnIdle bool                `json:"ask_if_human_present_on_idle"`
	OpenAIAccount            *AgentOpenAIAccount `json:"openai_account_connection"`
	RunDNCDetection          bool                `json:"run_do_not_call_detection"`
	LLMTemperature           float64             `json:"llm_temperature"`
}

func (a *Agent) UnmarshalJSON(data []byte) error {
	var id string
	if err := json.Unmarshal(data, &id); err == nil {
		a.ID = id
		return nil
	}

	type Alias Agent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

type AgentReqbase struct {
	Name                     string             `json:"name"`
	Prompt                   string             `json:"prompt"`
	Language                 Language           `json:"language"`
	Actions                  []string           `json:"actions"`
	Voice                    string             `json:"voice"`
	InitMsg                  string             `json:"initial_message,omitempty"`
	Webhook                  string             `json:"webhook,omitempty"`
	VectorDB                 string             `json:"vector_database,omitempty"`
	InterruptSense           InterruptSenseType `json:"interrupt_sensitivity"`
	CtxEndpint               string             `json:"context_endpoint,omitempty"`
	NoiseSuppression         bool               `json:"noise_suppression"`
	EndpointSense            EndpointSenseType  `json:"endpointing_sensitivity"`
	IVRNavMode               IVRNavModeType     `json:"ivr_navigation_mode"`
	Speed                    float32            `json:"conversation_speed"`
	InitMsgDelay             float64            `json:"initial_message_delay"`
	OpenAIModelOverride      string             `json:"openai_model_name_override,omitempty"`
	AsktIfHumanPresentOnIdle bool               `json:"ask_if_human_present_on_idle"`
	OpenAIAccount            *OpenAIAccount     `json:"openai_account_connection"`
	RunDNCDetection          bool               `json:"run_do_not_call_detection"`
	LLMTemperature           float64            `json:"llm_temperature"`
}

type CreateAgentReq struct {
	AgentReqbase
}

func (a CreateAgentReq) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.AgentReqbase)
}

type UpdateAgentReq struct {
	AgentReqbase
}

func (a UpdateAgentReq) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.AgentReqbase)
}

func (c *Client) ListAgents(ctx context.Context, paging *PageParams) (*Agents, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/agents/list")
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

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		actions := new(Agents)
		if err := json.NewDecoder(resp.Body).Decode(actions); err != nil {
			return nil, err
		}
		return actions, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
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

func (c *Client) GetAgent(ctx context.Context, agentID string) (*Agent, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/agents")
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
	q.Add("id", agentID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Agent)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
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

func (c *Client) CreateAgent(ctx context.Context, createReq *CreateAgentReq) (*Agent, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/agents/create")
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

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Agent)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
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

func (c *Client) UpdateAgent(ctx context.Context, actionID string, updateReq *UpdateAgentReq) (*Agent, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/agents/update")
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

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Agent)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
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
