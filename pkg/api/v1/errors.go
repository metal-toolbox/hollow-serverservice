package hollow

import (
	"fmt"
	"net/http"
)

// ClientError is returned when invalid arguments are provided to the Narwhal client
type ClientError struct {
	Message string
}

// Error returns the ClientError in string format
func (e *ClientError) Error() string {
	return fmt.Sprintf("hollow client error: %s", e.Message)
}

func newClientError(msg string) *ClientError {
	return &ClientError{
		Message: msg,
	}
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusMultiStatus {
		return newClientError(fmt.Sprintf("invalid response code: %d", resp.StatusCode))
	}

	return nil
}
