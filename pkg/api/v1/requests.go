package serverservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.hollow.sh/toolbox/version"
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
	return fmt.Sprintf("go-hollow-client (%s)", version.String())
}

func (c *Client) do(req *http.Request, result interface{}) error {
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.authToken))
	req.Header.Set("User-Agent", userAgentString())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if err := ensureValidServerResponse(resp); err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}
