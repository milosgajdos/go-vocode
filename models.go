package vocode

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
