package vocode

import (
	"os"

	"github.com/milosgajdos/go-vocode/client"
)

const (
	// BaseURL is OpenAI HTTP API base URL.
	BaseURL = "https://api.vocode.dev"
	// APIV2 V1 version.
	APIV1 = "v1"
)

// Client is an OpenAI HTTP API client.
type Client struct {
	opts Options
}

type Options struct {
	APIKey     string
	UserID     string
	BaseURL    string
	Version    string
	HTTPClient *client.HTTP
}

// Option is functional graph option.
type Option func(*Options)

// NewClient creates a new HTTP API client and returns it.
// By default it reads the secret key from PLAYHT_SECRET_KEY env var
// and user ID from PLAYHT_USER_ID env var and uses
// the default http client for making the HTTP api requests.
func NewClient(opts ...Option) *Client {
	options := Options{
		APIKey:     os.Getenv("VOCODE_API_KEY"),
		BaseURL:    BaseURL,
		Version:    APIV1,
		HTTPClient: client.NewHTTP(),
	}

	for _, apply := range opts {
		apply(&options)
	}

	return &Client{
		opts: options,
	}
}

// WithAPIKey sets the secret key.
func WithAPIKey(apiKey string) Option {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

// WithBaseURL sets the API base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

// WithVersion sets the API version.
func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(httpClient *client.HTTP) Option {
	return func(o *Options) {
		o.HTTPClient = httpClient
	}
}
