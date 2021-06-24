package hollow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var apiVersion = "v1"

// Client has the ability to talk to a hollow server running at the given URI
type Client struct {
	url        string
	authToken  string
	BIOSConfig BIOSConfigService
	Hardware   HardwareService
}

// NewClient will initialize a new hollow client with the given auth token and URL
func NewClient(authToken, url string) (*Client, error) {
	if authToken == "" {
		return nil, newClientError("failed to initialize: no auth token provided")
	}

	if url == "" {
		return nil, newClientError("failed to initialize: no hollow api url provided")
	}

	c := &Client{
		url:       url,
		authToken: authToken,
	}

	c.BIOSConfig = &BIOSConfigServiceClient{client: c}
	c.Hardware = &HardwareServiceClient{client: c}

	return c, nil
}

func newGetRequest(ctx context.Context, uri, path string) (*http.Request, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/api/%s/%s", uri, apiVersion, path))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func newPostRequest(ctx context.Context, uri, path string, body interface{}) (*http.Request, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/api/%s/%s", uri, apiVersion, path))
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}
