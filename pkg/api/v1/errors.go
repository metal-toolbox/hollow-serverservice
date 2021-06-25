package hollow

import (
	"fmt"
	"io"
	"net/http"
)

// ClientError is returned when invalid arguments are provided to the client
type ClientError struct {
	Message string
}

// ServerError is returned when the client receives an error back from the server
type ServerError struct {
	Message    string
	StatusCode int
}

// Error returns the ClientError in string format
func (e *ClientError) Error() string {
	return fmt.Sprintf("hollow client error: %s", e.Message)
}

// Error returns the ServerError in string format
func (e *ServerError) Error() string {
	return fmt.Sprintf("hollow client received a server error: status_code: %d, message: %s", e.StatusCode, e.Message)
}

func newClientError(msg string) *ClientError {
	return &ClientError{
		Message: msg,
	}
}

func ensureValidServerResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusMultiStatus {
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			b = []byte("failed to read body")
		}

		return &ServerError{
			StatusCode: resp.StatusCode,
			Message:    string(b),
		}
	}

	return nil
}
