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

type TelProvider string

const (
	Vonage TelProvider = "vonage"
	Twilio TelProvider = "twilio"
)

type NumbersItem struct {
	ID             string         `json:"id"`
	UserID         string         `json:"user_id"`
	Label          string         `json:"label"`
	Number         string         `json:"number"`
	TelAccountID   string         `json:"telephony_account_connection"`
	TelProvider    TelProvider    `json:"telephony_provider"`
	InboundAgentID string         `json:"inbound_agent,omitempty"`
	OutboundOnly   bool           `json:"outbound_only"`
	Active         bool           `json:"active"`
	ExampleCtx     map[string]any `json:"example_context,omitempty"`
}

type Numbers struct {
	Items []NumbersItem `json:"items"`
	*Paging
}

type Field struct {
	Type  string
	Label string
	Name  string
	Desc  string
}

type Template struct {
	ID         string   `json:"id"`
	UserID     string   `json:"user_id"`
	Label      string   `json:"label"`
	ReqCtxKeys []string `json:"required_context_keys"`
}

type Prompt struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Content     string    `json:"content"`
	Fields      []Field   `json:"collect_fields"`
	CtxEndpoint string    `json:"context_endpoint"`
	Template    *Template `json:"prompt_template"`
}

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
	ID      string                    `json:"id"`
	UserID  string                    `json:"user_id"`
	Type    ActionType                `json:"type"`
	Config  *TransferCallActionConfig `json:"config"`
	Trigger interface{}               `json:"action_trigger"`
}

type EndConversationAction struct {
	ID      string      `json:"id"`
	UserID  string      `json:"user_id"`
	Type    ActionType  `json:"type"`
	Config  map[any]any `json:"config"`
	Trigger interface{} `json:"action_trigger"`
}

type DTMFAction struct {
	ID      string      `json:"id"`
	UserID  string      `json:"user_id"`
	Type    ActionType  `json:"type"`
	Config  map[any]any `json:"config"`
	Trigger interface{} `json:"action_trigger"`
}

type AddToConfConfig struct {
	PhoneNr            string `json:"phone_number"`
	PlacePrimaryOnHold bool   `json:"place_primary_on_hold"`
}

type AddToConfAction struct {
	ID      string           `json:"id"`
	UserID  string           `json:"user_id"`
	Type    ActionType       `json:"type"`
	Config  *AddToConfConfig `json:"config"`
	Trigger interface{}      `json:"action_trigger"`
}

type SetHoldAction struct {
	ID      string      `json:"id"`
	UserID  string      `json:"user_id"`
	Type    ActionType  `json:"type"`
	Config  map[any]any `json:"config"`
	Trigger interface{} `json:"action_trigger"`
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
	ID     string                `json:"id"`
	UserID string                `json:"user_id"`
	Type   ActionType            `json:"type"`
	Config *ExternalActionConfig `json:"config"`
}

type VoiceType string

const (
	Azure      VoiceType = "voice_azure"
	Rime       VoiceType = "voice_rime"
	ElevenLabs VoiceType = "voice_eleven_labs"
	PlayHt     VoiceType = "voice_play_ht"
)

type VoiceModel string

const (
	Mist VoiceModel = "mist"
	V1   VoiceModel = "v1"
)

type AzureVoice struct {
	Type  VoiceType `json:"type"`
	Name  string    `json:"voice_name"`
	Pitch int       `json:"pitch"`
	Rate  int       `json:"rate"`
}

type RimeVoice struct {
	Type       VoiceType  `json:"type"`
	Speaker    string     `json:"speaker"`
	SpeedAlpha string     `json:"speed_alpha"`
	ModelID    VoiceModel `json:"model_id"`
}

type ElevenLabsVoice struct {
	Type           VoiceType `json:"type"`
	APIKey         string    `json:"api_key"`
	ModelID        string    `json:"model_id"`
	VoiceID        string    `json:"voice_id"`
	Stability      int       `json:"stability"`
	SimBoost       int       `json:"similarity_boost"`
	OptimStream    int       `json:"optimize_streaming_latency"`
	ExpInputStream bool      `json:"experimental_input_streaming"`
}

type PlayHtVersion string

const (
	PlayHtV1 PlayHtVersion = "1"
	PlayHtV2 PlayHtVersion = "2"
)

type PlayHtQuality string

const (
	Faster  PlayHtQuality = "faster"
	Draft   PlayHtQuality = "draft"
	Low     PlayHtQuality = "low"
	Medium  PlayHtQuality = "medium"
	High    PlayHtQuality = "high"
	Premium PlayHtQuality = "premium"
)

type PlayHtVoice struct {
	Type             VoiceType     `json:"type"`
	VoiceID          string        `json:"voice_id"`
	APIUserID        string        `json:"api_user_id"`
	APIKey           string        `json:"api_key"`
	Version          PlayHtVersion `json:"version"`
	Speed            int           `json:"speed"`
	Quality          PlayHtQuality `json:"quality"`
	Temp             int           `json:"temperature"`
	TopP             int           `json:"top_p"`
	TextGuidance     string        `json:"text_guidance"`
	VoiceGuidance    string        `json:"voice_guidance"`
	ExpRemoveSilence bool          `json:"experimental_remove_silence"`
}

// TODO: UnmarhalJSON
type Voice struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	VoiceType any    `json:"-"`
}

type Event string

const (
	EventMessage          Event = "event_message"
	EventAction           Event = "event_action"
	EventCallConnected    Event = "event_phone_call_connected"
	EventCallEnded        Event = "event_phone_call_ended"
	EventCallDidntConnect Event = "event_phone_call_did_not_connect"
	EventTranscript       Event = "event_transcript"
	EventRecording        Event = "event_recording"
	EventHumanDetection   Event = "event_human_detection"
)

type WebhookMethod string

const (
	Get  WebhookMethod = "GET"
	Post WebhookMethod = "POST"
)

type Webhook struct {
	ID     string        `json:"id"`
	UserID string        `json:"user_id"`
	Subs   []Event       `json:"subscriptions"`
	URL    string        `json:"url"`
	Method WebhookMethod `json:"method"`
}

type VectorDBType string

const (
	PineCone VectorDBType = "vector_database_pinecone"
)

type VectorDB struct {
	ID     string       `json:"id"`
	UserID string       `json:"user_id"`
	Type   VectorDBType `json:"type"`
	Index  string       `json:"index"`
	APIKey string       `json:"api_key"`
	APIEnv string       `json:"api_environment"`
}

type OpenAICreds struct {
	APIKey string `json:"openai_api_key"`
}

type AcctConnectionType string

const (
	OpenaiConnType AcctConnectionType = "account_connection_openai"
	TwilioConnType AcctConnectionType = "account_connection_twilio"
)

type OpenAIAccount struct {
	ID     string             `json:"id"`
	UserID string             `json:"user_id"`
	Type   AcctConnectionType `json:"type"`
	Creds  *OpenAICreds       `json:"credentials"`
}

type InterruptSenseType string

const (
	LowInterrupt  InterruptSenseType = "low"
	HighInterrupt InterruptSenseType = "high"
)

type EndpointSenseType string

const (
	AutoEndpoint      EndpointSenseType = "auto"
	RelaxedEndpoint   EndpointSenseType = "relaxed"
	SensitiveEndpoint EndpointSenseType = "sensitive"
)

type IVRNavModeType string

const (
	DefaultIVRMode IVRNavModeType = "default"
	OffIVRMode     IVRNavModeType = "off"
)

type InboundAgent struct {
	ID                  string             `json:"id"`
	UserID              string             `json:"user_id"`
	Name                string             `json:"name"`
	Prompt              *Prompt            `json:"prompt"`
	Language            Language           `json:"language"`
	Actions             []interface{}      `json:"actions"`
	Voice               *Voice             `json:"voice"`
	InitMsg             string             `json:"initial_msg"`
	Webhook             *Webhook           `json:"webhook"`
	VectorDB            *VectorDB          `json:"vector_database"`
	InterruptSense      InterruptSenseType `json:"interrupt_sensitivity"`
	CtxEndpint          string             `json:"context_endpint"`
	NoiseSuppression    bool               `json:"noise_suppression"`
	EndpointSense       EndpointSenseType  `json:"endpointing_sensitivity"`
	IVRNavMode          IVRNavModeType     `json:"ivr_navigation_mode"`
	Speed               int                `json:"conversation_speed"`
	InitMsgDelay        int                `json:"initial_message_delay"`
	OpenAIModelOverride bool               `json:"openai_model_name_override"`
	AsktIfHumanPresent  bool               `json:"ask_if_human_present_on_idle"`
	OpenAIAccount       *OpenAIAccount     `json:"openai_account_connection"`
	RunDNCDetecion      bool               `json:"run_do_not_call_detection"`
	LLMTemperature      int                `json:"llm_temperature"`
}

type TelAccount struct {
	ID               string             `json:"id"`
	UserID           string             `json:"user_id"`
	Type             AcctConnectionType `json:"type"`
	Credentials      map[string]any     `json:"credentials"`
	SteeringPool     []string           `json:"steering_pool"`
	SupportAnyCaller bool               `json:"account_supports_any_caller_id"`
}

type Number struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id"`
	Active       bool           `json:"active"`
	Label        string         `json:"label"`
	InboundAgent *InboundAgent  `json:"inbound_agent"`
	OutboundOnly bool           `json:"outbound_only"`
	ExampleCtx   map[string]any `json:"example_context"`
	Number       string         `json:"number"`
	TelProvider  TelProvider    `json:"telephony_provider"`
	TelAccount   *TelAccount    `json:"telephony_account_connection"`
}

type BuyNumberReq struct {
	AreaCode     string      `json:"area_code"`
	TelProvider  TelProvider `json:"telephony_provider"`
	TelAccountID string      `json:"telephony_account_connection"`
}

type UpdateNumberReq struct {
	Label          string         `json:"label"`
	InboundAgentID string         `json:"inbound_agent"`
	OutboundOnly   bool           `json:"outbound_only"`
	ExampleCtx     map[string]any `json:"example_context"`
}

func (c *Client) ListNumbers(ctx context.Context, paging *PageParams) (*Numbers, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/numbers/list")
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
		nrResp := new(Numbers)
		if err := json.NewDecoder(resp.Body).Decode(nrResp); err != nil {
			return nil, err
		}
		return nrResp, nil
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

func (c *Client) GetNumber(ctx context.Context, number string) (*Number, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/numbers")
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
	q.Add("phone_number", number)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		nrResp := new(Number)
		if err := json.NewDecoder(resp.Body).Decode(nrResp); err != nil {
			return nil, err
		}
		return nrResp, nil
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

func (c *Client) BuyNumber(ctx context.Context, buyReq *BuyNumberReq) (*Number, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/numbers/buy")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(buyReq); err != nil {
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
		nrResp := new(Number)
		if err := json.NewDecoder(resp.Body).Decode(nrResp); err != nil {
			return nil, err
		}
		return nrResp, nil
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

func (c *Client) UpdateNumber(ctx context.Context, number string, updateReq *UpdateNumberReq) (*Number, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/numbers/update")
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
	q.Add("phone_number", number)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		nrResp := new(Number)
		if err := json.NewDecoder(resp.Body).Decode(nrResp); err != nil {
			return nil, err
		}
		return nrResp, nil
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

func (c *Client) CancelNumber(ctx context.Context, number string) (*Number, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/numbers/cancel")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("phone_number", number)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		nrResp := new(Number)
		if err := json.NewDecoder(resp.Body).Decode(nrResp); err != nil {
			return nil, err
		}
		return nrResp, nil
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
