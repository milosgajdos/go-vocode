package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/milosgajdos/go-vocode/client"
)

const (
	// UserAgent is default User-Agent header value.
	UserAgent = "github.com/milosgajdos/go-vocode"
)

// NewHTTP creates a new HTTP request from the provided parameters and returns it.
// If the passed in context is nil, it creates a new background context.
// If the provided body is nil, it gets initialized to bytes.Reader.
// If no Content-Type has been set via options it defaults to application/json.
func NewHTTP(ctx context.Context, method, url string, body io.Reader, opts ...HTTPOption) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if body == nil {
		body = &bytes.Reader{}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	for _, setOption := range opts {
		setOption(req)
	}

	// if no content-type is specified we default to json
	if ct := req.Header.Get("Content-Type"); len(ct) == 0 {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("User-Agent", UserAgent)

	return req, nil
}

// Do sends the HTTP request req using the client and returns the response.
func Do[T error](client *client.HTTP, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		return resp, nil
	}
	defer resp.Body.Close()

	var apiErr T
	if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
		return nil, jsonErr
	}

	return nil, apiErr
}

// PageParams contain key-value pairs with request paging parameters
type PageParams map[string]string

// HTTPOption is a HTTP request functional option.
type HTTPOption func(*http.Request)

// WithSetHeader sets the header key to value val.
func WithSetHeader(key, val string) HTTPOption {
	return func(req *http.Request) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Set(key, val)
	}
}

// WithAddHeader adds the val to key header.
func WithAddHeader(key, val string) HTTPOption {
	return func(req *http.Request) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Add(key, val)
	}
}

// WithBearer sets the Authorization header to the provided Bearer token.
func WithBearer(token string) HTTPOption {
	return func(req *http.Request) {
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// WithPageParams encodes paging params
func WithPageParams(params PageParams) HTTPOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}
