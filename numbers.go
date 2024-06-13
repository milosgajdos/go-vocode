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
	VonageTelProvider TelProvider = "vonage"
	TwilioTelProvider TelProvider = "twilio"
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
		numbers := new(Numbers)
		if err := json.NewDecoder(resp.Body).Decode(numbers); err != nil {
			return nil, err
		}
		return numbers, nil
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
		number := new(Number)
		if err := json.NewDecoder(resp.Body).Decode(number); err != nil {
			return nil, err
		}
		return number, nil
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
