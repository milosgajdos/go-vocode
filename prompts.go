package vocode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

type FieldType string

const (
	EmailFieldType FieldType = "field_type_email"
)

type Field struct {
	Type  FieldType `json:"field_type"`
	Label string    `json:"label"`
	Name  string    `json:"name"`
	Desc  string    `json:"description"`
}

type Template struct {
	ID         string   `json:"id"`
	UserID     string   `json:"user_id"`
	Label      string   `json:"label"`
	ReqCtxKeys []string `json:"required_context_keys"`
}

type Prompts struct {
	Items []Prompt `json:"items"`
	*Paging
}

type Prompt struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Content     string    `json:"content"`
	Fields      []Field   `json:"collect_fields"`
	CtxEndpoint string    `json:"context_endpoint"`
	Template    *Template `json:"prompt_template"`
}

func (p *Prompt) UnmarshalJSON(data []byte) error {
	var id string
	if err := json.Unmarshal(data, &id); err == nil {
		p.ID = id
		return nil
	}

	type Alias Prompt
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

type PromptReq struct {
	Content     string  `json:"content"`
	Fields      []Field `json:"collect_fields"`
	CtxEndpoint string  `json:"context_endpoint,omitempty"`
	Template    string  `json:"prompt_template,omitempty"`
}

type CreatePromptReq struct {
	PromptReq
}

type UpdatePromptReq struct {
	PromptReq
}

func (c *Client) ListPrompts(ctx context.Context, paging *PageParams) (*Prompts, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/prompts/list")
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

	prompts := new(Prompts)
	if err := json.NewDecoder(resp.Body).Decode(prompts); err != nil {
		return nil, err
	}
	return prompts, nil
}

func (c *Client) GetPrompt(ctx context.Context, id string) (*Prompt, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/prompts")
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

	prompts := new(Prompt)
	if err := json.NewDecoder(resp.Body).Decode(prompts); err != nil {
		return nil, err
	}
	return prompts, nil
}

func (c *Client) CreatePrompt(ctx context.Context, createReq *CreatePromptReq) (*Prompt, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/prompts/create")
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

	prompt := new(Prompt)
	if err := json.NewDecoder(resp.Body).Decode(prompt); err != nil {
		return nil, err
	}
	return prompt, nil
}

func (c *Client) UpdatePrompt(ctx context.Context, id string, updateReq *UpdatePromptReq) (*Prompt, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/prompts/update")
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

	prompt := new(Prompt)
	if err := json.NewDecoder(resp.Body).Decode(prompt); err != nil {
		return nil, err
	}
	return prompt, nil
}
