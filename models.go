package vocode

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
