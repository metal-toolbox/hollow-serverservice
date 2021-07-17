package hollow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.metalkube.net/hollow/pkg/version"
)

func newGetRequest(ctx context.Context, uri, path string) (*http.Request, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/api/%s/%s", uri, apiVersion, path))
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
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

	return http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), buf)
}

func newPutRequest(ctx context.Context, uri, path string, body interface{}) (*http.Request, error) {
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

	return http.NewRequestWithContext(ctx, http.MethodPut, requestURL.String(), buf)
}

func newDeleteRequest(ctx context.Context, uri, path string) (*http.Request, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%s/api/%s/%s", uri, apiVersion, path))
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(ctx, http.MethodDelete, requestURL.String(), nil)
}

func userAgentString() string {
	return fmt.Sprintf("hollow/%s (%s)", version.Version(), version.String())
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.authToken))
	req.Header.Set("User-Agent", userAgentString())

	return c.httpClient.Do(req)
}
