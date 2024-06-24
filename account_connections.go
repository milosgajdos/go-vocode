package vocode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

type AccountConnType string

const (
	AccountConnOpenAI AccountConnType = "account_connection_openai"
	AccountConnTwilio AccountConnType = "account_connection_twilio"
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
	var id string
	if err := json.Unmarshal(data, &id); err == nil {
		ta.ID = id
		return nil
	}

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
	case AccountConnOpenAI:
		var openaiAccount OpenAIAccount
		if err := json.Unmarshal(data, &openaiAccount); err != nil {
			return err
		}
		a.OpenAIAccount = &openaiAccount
	case AccountConnTwilio:
		var twillioAccount TwilioAccount
		if err := json.Unmarshal(data, &twillioAccount); err != nil {
			return err
		}
		a.TwilioAccount = &twillioAccount
	}

	return nil
}

type AccountConnReq struct {
	Type          AccountConnType `json:"type"`
	TwilioAccount *TwilioAccount  `json:"-"`
	OpenAIAccount *OpenAIAccount  `json:"-"`
}

func (a AccountConnReq) MarshalJSON() ([]byte, error) {
	type Alias AccountConnReq

	switch a.Type {
	case AccountConnOpenAI:
		return json.Marshal(&struct {
			*Alias
			*OpenAIAccount
		}{
			Alias:         (*Alias)(&a),
			OpenAIAccount: a.OpenAIAccount,
		})
	case AccountConnTwilio:
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
	AccountConnReq
}

func (a CreateAccountConnReq) MarshalJSON() ([]byte, error) {
	return a.AccountConnReq.MarshalJSON()
}

type UpdateAccountConnReq struct {
	AccountConnReq
}

func (a UpdateAccountConnReq) MarshalJSON() ([]byte, error) {
	return a.AccountConnReq.MarshalJSON()
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

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	accountConns := new(AccountConns)
	if err := json.NewDecoder(resp.Body).Decode(accountConns); err != nil {
		return nil, err
	}
	return accountConns, nil
}

func (c *Client) GetAccountConn(ctx context.Context, id string) (*AccountConn, error) {
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
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	accountConn := new(AccountConn)
	if err := json.NewDecoder(resp.Body).Decode(accountConn); err != nil {
		return nil, err
	}
	return accountConn, nil
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

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	accountConn := new(AccountConn)
	if err := json.NewDecoder(resp.Body).Decode(accountConn); err != nil {
		return nil, err
	}
	return accountConn, nil
}

func (c *Client) UpdateAccountConn(ctx context.Context, id string, updateReq *UpdateAccountConnReq) (*AccountConn, error) {
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
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	accountConn := new(AccountConn)
	if err := json.NewDecoder(resp.Body).Decode(accountConn); err != nil {
		return nil, err
	}
	return accountConn, nil
}
