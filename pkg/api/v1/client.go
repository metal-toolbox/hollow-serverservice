package hollow

import (
	"context"
	"net/http"
	"net/url"
)

var apiVersion = "v1"

// Client has the ability to talk to a hollow server running at the given URI
type Client struct {
	url                 string
	authToken           string
	httpClient          Doer
	Server              ServerService
	ServerComponentType ServerComponentTypeService
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

	c.Server = &ServerServiceClient{client: c}
	c.ServerComponentType = &ServerComponentTypeServiceClient{client: c}

	return c, nil
}

// SetToken allows you to change the token of a client
func (c *Client) SetToken(token string) {
	c.authToken = token
}

// post provides a reusable method for a standard POST to a hollow server
func (c *Client) post(ctx context.Context, path string, body interface{}) (*ServerResponse, error) {
	request, err := newPostRequest(ctx, c.url, path, body)
	if err != nil {
		return nil, err
	}

	r := ServerResponse{}

	if err := c.do(request, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// put provides a reusable method for a standard PUT to a hollow server
func (c *Client) put(ctx context.Context, path string, body interface{}) (*ServerResponse, error) {
	request, err := newPutRequest(ctx, c.url, path, body)
	if err != nil {
		return nil, err
	}

	r := ServerResponse{}

	if err := c.do(request, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

type queryParams interface {
	setQuery(url.Values)
}

// list provides a reusable method for a standard list to a hollow server
func (c *Client) list(ctx context.Context, path string, params queryParams, resp interface{}) error {
	request, err := newGetRequest(ctx, c.url, path)
	if err != nil {
		return err
	}

	if params != nil {
		q := request.URL.Query()
		params.setQuery(q)
		request.URL.RawQuery = q.Encode()
	}

	return c.do(request, &resp)
}

// get provides a reusable method for a standard GET of a single item
func (c *Client) get(ctx context.Context, path string, resp interface{}) error {
	request, err := newGetRequest(ctx, c.url, path)
	if err != nil {
		return err
	}

	return c.do(request, &resp)
}

// post provides a reusable method for a standard post to a hollow server
func (c *Client) delete(ctx context.Context, path string) (*ServerResponse, error) {
	request, err := newDeleteRequest(ctx, c.url, path)
	if err != nil {
		return nil, err
	}

	var r ServerResponse

	return &r, c.do(request, &r)
}
