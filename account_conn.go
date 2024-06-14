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

type AccountConnType string

const (
	OpenaiConnType AccountConnType = "account_connection_openai"
	TwilioConnType AccountConnType = "account_connection_twilio"
)

type OpenAICreds struct {
	APIKey string `json:"openai_api_key"`
}

type OpenAIAccount struct {
	Creds *OpenAICreds `json:"credentials"`
}

type TwilioCreds struct {
	AccountID string `json:"twilio_account_sid"`
	AuthToken string `json:"twilio_auth_token"`
}

type TwilioAccount struct {
	Creds             *TwilioCreds `json:"credentials"`
	SteeringPool      []string     `json:"steering_pool"`
	SupportsAnyCaller bool         `json:"account_supports_any_caller_id"`
}

type AccountConnsBase struct {
	ID     string          `json:"id"`
	UserID string          `json:"user_id"`
	Type   AccountConnType `json:"type"`
}

type AccountConns struct {
	Items []AccountConn `json:"items"`
	*Paging
}

type AccountConn struct {
	AccountConnsBase
	TwilioAccount *TwilioAccount `json:",omitempty"`
	OpenAIAccount *OpenAIAccount `json:",omitempty"`
}

func (a *AccountConn) UnmarshalJSON(data []byte) error {
	var base AccountConnsBase
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}
	a.AccountConnsBase = base

	switch a.Type {
	case OpenaiConnType:
		var openaiAccount OpenAIAccount
		if err := json.Unmarshal(data, &openaiAccount); err != nil {
			return err
		}
		a.OpenAIAccount = &openaiAccount
	case TwilioConnType:
		var twillioAccount TwilioAccount
		if err := json.Unmarshal(data, &twillioAccount); err != nil {
			return err
		}
		a.TwilioAccount = &twillioAccount
	}

	return nil
}

type AccountConnReqBase struct {
	Type          AccountConnType `json:"type"`
	TwilioAccount *TwilioAccount  `json:"-"`
	OpenAIAccount *OpenAIAccount  `json:"-"`
}

func (a AccountConnReqBase) MarshalJSON() ([]byte, error) {
	type Alias AccountConnReqBase

	switch a.Type {
	case OpenaiConnType:
		return json.Marshal(&struct {
			*Alias
			*TwilioAccount
		}{
			Alias:         (*Alias)(&a),
			TwilioAccount: a.TwilioAccount,
		})
	case TwilioConnType:
		return json.Marshal(&struct {
			*Alias
			*OpenAIAccount
		}{
			Alias:         (*Alias)(&a),
			OpenAIAccount: a.OpenAIAccount,
		})
	default:
		return nil, fmt.Errorf("unsupported account connection type: %s", a.Type)
	}
}

type CreateAccountConnReq struct {
	AccountConnReqBase
}

func (a CreateAccountConnReq) MarshalJSON() ([]byte, error) {
	return a.AccountConnReqBase.MarshalJSON()
}

type UpdateAccountConnReq struct {
	AccountConnReqBase
}

func (a UpdateAccountConnReq) MarshalJSON() ([]byte, error) {
	return a.AccountConnReqBase.MarshalJSON()
}

func (c *Client) ListAccountConns(ctx context.Context, paging *PageParams) (*AccountConns, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/account_connections/list")
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
		actions := new(AccountConns)
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

func (c *Client) GetAccountConn(ctx context.Context, voiceID string) (*AccountConn, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/account_connections")
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
		action := new(AccountConn)
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

func (c *Client) CreateAccountConn(ctx context.Context, createReq *CreateAccountConnReq) (*AccountConn, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/account_connections/create")
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
		action := new(AccountConn)
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

func (c *Client) UpdateAccountConn(ctx context.Context, actionID string, updateReq *UpdateAccountConnReq) (*AccountConn, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/account_connections/update")
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
		action := new(AccountConn)
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
