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
	OpenAIConnType AccountConnType = "account_connection_openai"
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

type TelAccountConn struct {
	ID               string          `json:"id"`
	UserID           string          `json:"user_id"`
	Type             AccountConnType `json:"type"`
	Credentials      map[string]any  `json:"credentials"`
	SteeringPool     []string        `json:"steering_pool"`
	SupportAnyCaller bool            `json:"account_supports_any_caller_id"`
}

func (ta *TelAccountConn) UnmarshalJSON(data []byte) error {
	// Check if the data is a plain string ID
	var id string
	if err := json.Unmarshal(data, &id); err == nil {
		ta.ID = id
		return nil
	}

	// Otherwise, unmarshal as a full TelAccountConn object
	type Alias TelAccountConn
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(ta),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

type TellMetadataType string

const (
	TelMetadataVonage TellMetadataType = "telephony_metadata_vonage"
	TelMetadataTwilio TellMetadataType = "telephony_metadata_twilio"
)

type VonageTelMetadata struct {
	Type TellMetadataType `json:"type"`
}

type TwilioTelMetadata struct {
	Type               TellMetadataType `json:"type"`
	CallSID            string           `json:"call_sid"`
	CallStatus         string           `json:"call_status"`
	TransferCallSID    string           `json:"transfer_call_sid"`
	TransferCallStatus string           `json:"transfer_call_status"`
	ConferenceSID      string           `json:"conference_sid"`
}

type TelMetadataBase struct {
	Type TellMetadataType `json:"type"`
}

type TelMetadata struct {
	TelMetadataBase
	*VonageTelMetadata
	*TwilioTelMetadata
}

func (t *TelMetadata) UnmarshalJSON(data []byte) error {
	var base TelMetadataBase
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}
	t.TelMetadataBase = base

	switch t.TelMetadataBase.Type {
	case TelMetadataVonage:
		t.VonageTelMetadata = &VonageTelMetadata{
			Type: TelMetadataVonage,
		}
		return nil
	case TelMetadataTwilio:
		t.TwilioTelMetadata = &TwilioTelMetadata{}
		return json.Unmarshal(data, t.TwilioTelMetadata)
	}

	return nil
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
	case OpenAIConnType:
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
	case OpenAIConnType:
		return json.Marshal(&struct {
			*Alias
			*OpenAIAccount
		}{
			Alias:         (*Alias)(&a),
			OpenAIAccount: a.OpenAIAccount,
		})
	case TwilioConnType:
		return json.Marshal(&struct {
			*Alias
			*TwilioAccount
		}{
			Alias:         (*Alias)(&a),
			TwilioAccount: a.TwilioAccount,
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

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
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

func (c *Client) GetAccountConn(ctx context.Context, acctConnID string) (*AccountConn, error) {
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
	q.Add("id", acctConnID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
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

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
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

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
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
