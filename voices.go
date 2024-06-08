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

type VoiceType string

const (
	AzureVoiceType      VoiceType = "voice_azure"
	RimeVoiceType       VoiceType = "voice_rime"
	ElevenLabsVoiceType VoiceType = "voice_eleven_labs"
	PlayHtVoiceType     VoiceType = "voice_play_ht"
)

type RimeVoiceModel string

const (
	MistRimeVoiceModel   RimeVoiceModel = "mist"
	V1RimeRimeVoiceModel RimeVoiceModel = "v1"
)

type AzureVoice struct {
	Name  string `json:"voice_name"`
	Pitch int    `json:"pitch"`
	Rate  int    `json:"rate"`
}

type RimeVoice struct {
	Speaker    string         `json:"speaker"`
	SpeedAlpha float64        `json:"speed_alpha"`
	ModelID    RimeVoiceModel `json:"model_id"`
}

type ElevenLabsVoice struct {
	APIKey         string `json:"api_key"`
	ModelID        string `json:"model_id"`
	VoiceID        string `json:"voice_id"`
	Stability      int    `json:"stability"`
	SimBoost       int    `json:"similarity_boost"`
	OptimStream    int    `json:"optimize_streaming_latency"`
	ExpInputStream bool   `json:"experimental_input_streaming"`
}

type PlayHtVersion string

const (
	PlayHtV1 PlayHtVersion = "1"
	PlayHtV2 PlayHtVersion = "2"
)

type PlayHtQuality string

const (
	FasterPlayHtQuality  PlayHtQuality = "faster"
	DraftPlayHtQuality   PlayHtQuality = "draft"
	LowPlayHtQuality     PlayHtQuality = "low"
	MediumPlayHtQuality  PlayHtQuality = "medium"
	HighPlayHtQuality    PlayHtQuality = "high"
	PremiumPlayHtQuality PlayHtQuality = "premium"
)

type PlayHtVoice struct {
	VoiceID          string        `json:"voice_id"`
	APIUserID        string        `json:"api_user_id"`
	APIKey           string        `json:"api_key"`
	Version          PlayHtVersion `json:"version"`
	Quality          PlayHtQuality `json:"quality"`
	Speed            float32       `json:"speed"`
	Temp             float32       `json:"temperature"`
	TopP             int           `json:"top_p"`
	TextGuidance     string        `json:"text_guidance"`
	VoiceGuidance    string        `json:"voice_guidance"`
	ExpRemoveSilence bool          `json:"experimental_remove_silence"`
}

type VoiceBase struct {
	ID     string    `json:"id"`
	UserID string    `json:"user_id"`
	Type   VoiceType `json:"type"`
}

type Voice struct {
	VoiceBase
	AzureVoice      *AzureVoice      `json:",omitempty"`
	RimeVoice       *RimeVoice       `json:",omitempty"`
	ElevenLabsVoice *ElevenLabsVoice `json:",omitempty"`
	PlayHtVoice     *PlayHtVoice     `json:",omitempty"`
}

type Voices struct {
	Items []Voice `json:"items"`
	*Paging
}

type VoiceReqBase struct {
	Type            VoiceType        `json:"type"`
	AzureVoice      *AzureVoice      `json:"-"`
	RimeVoice       *RimeVoice       `json:"-"`
	ElevenLabsVoice *ElevenLabsVoice `json:"-"`
	PlayHtVoice     *PlayHtVoice     `json:"-"`
}

func (v VoiceReqBase) MarshalJSON() ([]byte, error) {
	type Alias VoiceReqBase

	switch v.Type {
	case AzureVoiceType:
		return json.Marshal(&struct {
			*Alias
			*AzureVoice
		}{
			Alias:      (*Alias)(&v),
			AzureVoice: v.AzureVoice,
		})
	case RimeVoiceType:
		return json.Marshal(&struct {
			*Alias
			*RimeVoice
		}{
			Alias:     (*Alias)(&v),
			RimeVoice: v.RimeVoice,
		})
	case ElevenLabsVoiceType:
		return json.Marshal(&struct {
			*Alias
			*ElevenLabsVoice
		}{
			Alias:           (*Alias)(&v),
			ElevenLabsVoice: v.ElevenLabsVoice,
		})
	case PlayHtVoiceType:
		return json.Marshal(&struct {
			*Alias
			*PlayHtVoice
		}{
			Alias:       (*Alias)(&v),
			PlayHtVoice: v.PlayHtVoice,
		})
	default:
		return nil, fmt.Errorf("unsupported voice type: %s", v.Type)
	}
}

type CreateVoiceReq struct {
	VoiceReqBase
}

func (v CreateVoiceReq) MarshalJSON() ([]byte, error) {
	return v.VoiceReqBase.MarshalJSON()
}

type UpdateVoiceReq struct {
	VoiceReqBase
}

func (v UpdateVoiceReq) MarshalJSON() ([]byte, error) {
	return v.VoiceReqBase.MarshalJSON()
}

func (v *Voice) UnmarshalJSON(data []byte) error {
	var base VoiceBase
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}
	v.VoiceBase = base

	switch v.Type {
	case AzureVoiceType:
		var azureVoice AzureVoice
		if err := json.Unmarshal(data, &azureVoice); err != nil {
			return err
		}
		v.AzureVoice = &azureVoice

	case RimeVoiceType:
		var rimeVoice RimeVoice
		if err := json.Unmarshal(data, &rimeVoice); err != nil {
			return err
		}
		v.RimeVoice = &rimeVoice

	case ElevenLabsVoiceType:
		var elevenLabsVoice ElevenLabsVoice
		if err := json.Unmarshal(data, &elevenLabsVoice); err != nil {
			return err
		}
		v.ElevenLabsVoice = &elevenLabsVoice

	case PlayHtVoiceType:
		var playHtVoice PlayHtVoice
		if err := json.Unmarshal(data, &playHtVoice); err != nil {
			return err
		}
		v.PlayHtVoice = &playHtVoice
	}

	return nil
}

func (c *Client) ListVoices(ctx context.Context, paging *PageParams) (*Voices, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/voices/list")
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
		actions := new(Voices)
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

func (c *Client) GetVoice(ctx context.Context, voiceID string) (*Voice, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/voices")
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
	q.Add("id", voiceID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Voice)
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

func (c *Client) CreateVoice(ctx context.Context, createReq *CreateVoiceReq) (*Voice, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/voices/create")
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
		action := new(Voice)
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

func (c *Client) UpdateVoice(ctx context.Context, actionID string, updateReq *UpdateVoiceReq) (*Voice, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/voices/update")
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
		action := new(Voice)
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
