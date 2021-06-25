package hollow

import (
	"net/http"
)

var apiVersion = "v1"

// Client has the ability to talk to a hollow server running at the given URI
type Client struct {
	url        string
	authToken  string
	httpClient Doer
	BIOSConfig BIOSConfigService
	Hardware   HardwareService
}

// Doer is an interface for an HTTP client that can make requests
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// NewClient will initialize a new hollow client with the given auth token and URL
func NewClient(authToken, url string, doerClient Doer) (*Client, error) {
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

	c.httpClient = doerClient
	if c.httpClient == nil {
		// Use the default client as a fallback if one isn't passed
		c.httpClient = http.DefaultClient
	}

	c.BIOSConfig = &BIOSConfigServiceClient{client: c}
	c.Hardware = &HardwareServiceClient{client: c}

	return c, nil
}
