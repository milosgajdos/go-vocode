package vocode

type NumbersItem struct {
	ID                 string         `json:"id"`
	UserID             string         `json:"user_id"`
	Active             bool           `json:"active"`
	Label              string         `json:"label"`
	InboundAgent       string         `json:"inbound_agent"`
	OutboundOnly       bool           `json:"outbound_only"`
	ExampleCtx         map[string]any `json:"example_context"`
	Number             string         `json:"number"`
	TeleProvider       TeleProvider   `json:"telephony_provider"`
	TeleAcctConnection string         `json:"telephony_account_connection"`
}

type Numbers struct {
	Items []NumbersItem `json:"items"`
	Pages *Pages        `json:"pages,omitempty"`
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

type TeleAccount struct {
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
	TeleProvider TeleProvider   `json:"telephony_provider"`
	TeleAccount  *TeleAccount   `json:"telephony_account_connection"`
}

type BuyNumberResp struct{}

type UpdateNumberResp struct{}

type CancelNumberResp struct{}

func (c *Client) GetNumbers() (*Numbers, error) {
	return nil, nil
}

func (c *Client) GetNumber() (*Number, error) {
	return nil, nil
}

func (c *Client) BuyNumber() (*BuyNumberResp, error) {
	return nil, nil
}

func (c *Client) UpdateNumber() (*UpdateNumberResp, error) {
	return nil, nil
}

func (c *Client) CancelNumber() (*CancelNumberResp, error) {
	return nil, nil
}

//func (c *Client) LinkNumber() error {
//	return nil, nil
//}
