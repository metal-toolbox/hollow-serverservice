package hollow_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

// MockHTTPRequestDoer implements the standard http.Client interface.
type MockHTTPRequestDoer struct {
	Response *http.Response
	Error    error
}

// Do mocks a HTTP request and response for use in testing the client without a server
func (md *MockHTTPRequestDoer) Do(req *http.Request) (*http.Response, error) {
	// For tests to make sure context is passed through correctly
	_, ok := req.Context().Deadline()
	if ok {
		return md.Response, context.DeadlineExceeded
	}

	// Add to response for test helping
	md.Response.Request = req

	// Make sure this isn't null to prevent null pointer errors in tests
	if md.Response.Body == nil {
		md.Response.Body = ioutil.NopCloser(bytes.NewBufferString("Hello World"))
	}

	return md.Response, md.Error
}

// mockClient that can be used for testing
func mockClient(body string, status int) *hollow.Client {
	mockDoer := &MockHTTPRequestDoer{
		Response: &http.Response{
			StatusCode: status,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		},
		Error: nil,
	}

	c, err := hollow.NewClient("mocked", "mocked", mockDoer)
	if err != nil {
		return nil
	}

	return c
}
