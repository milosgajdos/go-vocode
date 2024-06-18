package request

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestNewHTTPRequest(t *testing.T) {
	t.Parallel()
	t.Run("nil_body", func(t *testing.T) {
		t.Parallel()
		req, err := NewHTTP(context.TODO(), http.MethodGet, "http://foo.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		if req.Body == nil {
			t.Fatal("expected non-nil request body")
		}
	})
	t.Run("with_options", func(t *testing.T) {
		t.Parallel()
		token := "secret"
		options := []HTTPOption{
			WithBearer(token),
		}
		req, err := NewHTTP(context.TODO(), http.MethodGet, "http://foo.com", &bytes.Reader{}, options...)
		if err != nil {
			t.Fatal(err)
		}
		if req.Body == nil {
			t.Fatal("expected non-nil request body")
		}

		// check all default headers are set as well as the bearer one
		header := make(http.Header)
		header.Set("Content-Type", "application/json; charset=utf-8")
		header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		// NOTE: this header is set by default
		// on every request we create via NewHTTP.
		header.Set("User-Agent", UserAgent)
		if !reflect.DeepEqual(header, req.Header) {
			t.Fatalf("expected header: %+v, got: %+v", header, req.Header)
		}
	})
}

func TestHTTPReqOption(t *testing.T) {
	t.Parallel()
	t.Run("set_bearer", func(t *testing.T) {
		t.Parallel()
		req := &http.Request{}
		token := "token"
		WithBearer(token)(req)

		if authzVal := req.Header.Get("Authorization"); authzVal != fmt.Sprintf("Bearer %s", token) {
			t.Fatalf("expected Authorization header val: %+v, got: %+v", fmt.Sprintf("Bearer %s", token), authzVal)
		}
	})
	t.Run("set_header", func(t *testing.T) {
		t.Parallel()
		req := &http.Request{}
		key, val := "foo", "bar"
		WithSetHeader(key, val)(req)
		if headerVal := req.Header.Get(key); headerVal != val {
			t.Fatalf("expected header val: %+v, got: %+v", val, headerVal)
		}
	})

	t.Run("add_header", func(t *testing.T) {
		t.Parallel()
		key, val := "foo", "bar"
		req := &http.Request{
			Header: make(http.Header),
		}
		req.Header.Add(key, val)
		WithAddHeader(key, val)(req)

		if headerVals := req.Header.Values(key); !reflect.DeepEqual(headerVals, []string{val, val}) {
			t.Fatalf("expected header values: %+v, got: %+v", []string{val, val}, headerVals)
		}
	})
}
