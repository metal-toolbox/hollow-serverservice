package hollow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

var apiVersion = "v1"

// Client has the ability to talk to a hollow server running at the given URI
type Client struct {
	url                   string
	authToken             string
	httpClient            Doer
	Hardware              HardwareService
	HardwareComponentType HardwareComponentTypeService
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

	c.Hardware = &HardwareServiceClient{client: c}
	c.HardwareComponentType = &HardwareComponentTypeServiceClient{client: c}

	return c, nil
}

// post provides a reusable method for a standard POST to a hollow server
func (c *Client) post(ctx context.Context, path string, body interface{}) (*uuid.UUID, error) {
	request, err := newPostRequest(ctx, c.url, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r serverResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r.UUID, nil
}

// put provides a reusable method for a standard PUT to a hollow server
func (c *Client) put(ctx context.Context, path string, body interface{}) (*uuid.UUID, error) {
	request, err := newPutRequest(ctx, c.url, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(request)
	if err != nil {
		return nil, err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r serverResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r.UUID, nil
}

type queryParams interface {
	setQuery(url.Values)
}

// list provides a reusable method for a standard list to a hollow server
func (c *Client) list(ctx context.Context, path string, params queryParams, results interface{}) error {
	request, err := newGetRequest(ctx, c.url, path)
	if err != nil {
		return err
	}

	if params != nil {
		q := request.URL.Query()
		params.setQuery(q)
		request.URL.RawQuery = q.Encode()
	}

	resp, err := c.do(request)
	if err != nil {
		return err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&results)
}

// get provides a reusable method for a standard GET of a single item
func (c *Client) get(ctx context.Context, path string, results interface{}) error {
	request, err := newGetRequest(ctx, c.url, path)
	if err != nil {
		return err
	}

	resp, err := c.do(request)
	if err != nil {
		return err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&results)
}

// post provides a reusable method for a standard post to a hollow server
func (c *Client) delete(ctx context.Context, path string) error {
	request, err := newDeleteRequest(ctx, c.url, path)
	if err != nil {
		return err
	}

	resp, err := c.do(request)
	if err != nil {
		return err
	}

	return ensureValidServerResponse(resp)
}
