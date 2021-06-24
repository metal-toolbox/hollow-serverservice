package hollow

import (
	"fmt"
	"net/http"
)

// ClientError is returned when invalid arguments are provided to the client
type ClientError struct {
	Message string
}

// ServerError is returned when the client receives an error back from the server
type ServerError struct {
	Message string
}

// Error returns the ClientError in string format
func (e *ClientError) Error() string {
	return fmt.Sprintf("hollow client error: %s", e.Message)
}

// Error returns the ServerError in string format
func (e *ServerError) Error() string {
	return fmt.Sprintf("hollow client received a server error: %s", e.Message)
}

func newClientError(msg string) *ClientError {
	return &ClientError{
		Message: msg,
	}
}

func ensureValidServerResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusMultiStatus {
		return &ServerError{
			Message: fmt.Sprintf("invalid response code: %d", resp.StatusCode),
		}
	}

	return nil
}
