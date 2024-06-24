package vocode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

type VectorDBType string

const (
	PineConeVectorDB VectorDBType = "vector_database_pinecone"
)

type VectorDBs struct {
	Items []VectorDB `json:"items"`
	*Paging
}

type VectorDB struct {
	ID     string       `json:"id"`
	UserID string       `json:"user_id"`
	Type   VectorDBType `json:"type"`
	Index  string       `json:"index"`
	APIKey string       `json:"api_key"`
	APIEnv string       `json:"api_environment"`
}

func (v *VectorDB) UnmarshalJSON(data []byte) error {
	// Check if the data is a plain string ID
	var id string
	if err := json.Unmarshal(data, &id); err == nil {
		v.ID = id
		return nil
	}

	// Otherwise, unmarshal as a full TelAccountConn object
	type Alias VectorDB
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

type VectorDBReq struct {
	Type   VectorDBType `json:"type"`
	Index  string       `json:"index"`
	APIKey string       `json:"api_key"`
	APIEnv string       `json:"api_environment"`
}

type CreateVectorDBReq struct {
	VectorDBReq
}

type UpdateVectorDBReq struct {
	VectorDBReq
}

func (c *Client) ListVectorDBs(ctx context.Context, paging *PageParams) (*VectorDBs, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/vector_databases/list")
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

	vectorDBs := new(VectorDBs)
	if err := json.NewDecoder(resp.Body).Decode(vectorDBs); err != nil {
		return nil, err
	}
	return vectorDBs, nil
}

func (c *Client) GetVectorDB(ctx context.Context, id string) (*VectorDB, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/vector_databases")
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

	vectorDB := new(VectorDB)
	if err := json.NewDecoder(resp.Body).Decode(vectorDB); err != nil {
		return nil, err
	}
	return vectorDB, nil
}

func (c *Client) CreateVectorDB(ctx context.Context, createReq *CreateVectorDBReq) (*VectorDB, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/vector_databases/create")
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

	vectorDB := new(VectorDB)
	if err := json.NewDecoder(resp.Body).Decode(vectorDB); err != nil {
		return nil, err
	}
	return vectorDB, nil
}

func (c *Client) UpdateVectorDB(ctx context.Context, id string, updateReq *UpdateVectorDBReq) (*VectorDB, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/vector_databases/update")
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

	vectorDB := new(VectorDB)
	if err := json.NewDecoder(resp.Body).Decode(vectorDB); err != nil {
		return nil, err
	}
	return vectorDB, nil
}
