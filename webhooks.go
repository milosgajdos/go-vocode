package vocode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

type Event string

const (
	MessageEvent          Event = "event_message"
	ActionEvent           Event = "event_action"
	CallConnectedEvent    Event = "event_phone_call_connected"
	CallEndedEvent        Event = "event_phone_call_ended"
	CallDidntConnectEvent Event = "event_phone_call_did_not_connect"
	TranscriptEvent       Event = "event_transcript"
	RecordingEvent        Event = "event_recording"
	HumanDetectionEvent   Event = "event_human_detection"
)

type WebhookMethod string

const (
	GetWebhook  WebhookMethod = "GET"
	PostWebhook WebhookMethod = "POST"
)

type Webhooks struct {
	Items []Webhook `json:"items"`
	*Paging
}

type Webhook struct {
	ID     string        `json:"id,omitempty"`
	UserID string        `json:"user_id,omitempty"`
	Subs   []Event       `json:"subscriptions,omitempty"`
	URL    string        `json:"url,omitempty"`
	Method WebhookMethod `json:"method,omitempty"`
}

func (w *Webhook) UnmarshalJSON(data []byte) error {
	var id string
	if err := json.Unmarshal(data, &id); err == nil {
		w.ID = id
		return nil
	}

	type Alias Webhook
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(w),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

type WebhookReqBase struct {
	Subs   []Event       `json:"subscriptions"`
	URL    string        `json:"url"`
	Method WebhookMethod `json:"method"`
}

type CreateWebhookReq struct {
	WebhookReqBase
}

type UpdateWebhookReq struct {
	WebhookReqBase
}

func (c *Client) ListWebhooks(ctx context.Context, paging *PageParams) (*Webhooks, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/webhooks/list")
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

	webhooks := new(Webhooks)
	if err := json.NewDecoder(resp.Body).Decode(webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (c *Client) GetWebhook(ctx context.Context, id string) (*Webhook, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/webhooks")
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

	webhook := new(Webhook)
	if err := json.NewDecoder(resp.Body).Decode(webhook); err != nil {
		return nil, err
	}
	return webhook, nil
}

func (c *Client) CreateWebhook(ctx context.Context, createReq *CreateWebhookReq) (*Webhook, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/webhooks/create")
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

	webhook := new(Webhook)
	if err := json.NewDecoder(resp.Body).Decode(webhook); err != nil {
		return nil, err
	}
	return webhook, nil
}

func (c *Client) UpdateWebhook(ctx context.Context, id string, updateReq *UpdateWebhookReq) (*Webhook, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/webhooks/update")
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

	webhook := new(Webhook)
	if err := json.NewDecoder(resp.Body).Decode(webhook); err != nil {
		return nil, err
	}
	return webhook, nil
}
