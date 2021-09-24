package serverservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNoNextPage is the error returned when there is not an additional page of resources
	ErrNoNextPage = errors.New("no next page found")
)

// ClientError is returned when invalid arguments are provided to the client
type ClientError struct {
	Message string
}

// ServerError is returned when the client receives an error back from the server
type ServerError struct {
	Message      string `json:"message"`
	ErrorMessage string `json:"error"`
	StatusCode   int
}

// Error returns the ClientError in string format
func (e *ClientError) Error() string {
	return fmt.Sprintf("hollow client error: %s", e.Message)
}

// Error returns the ServerError in string format
func (e ServerError) Error() string {
	return fmt.Sprintf("hollow client received a server error - response code: %d, message: %s, details: %s", e.StatusCode, e.Message, e.ErrorMessage)
}

func newClientError(msg string) *ClientError {
	return &ClientError{
		Message: msg,
	}
}

func ensureValidServerResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusMultiStatus {
		defer resp.Body.Close()

		var se ServerError

		se.StatusCode = resp.StatusCode

		if err := json.NewDecoder(resp.Body).Decode(&se); err != nil {
			se.ErrorMessage = "failed to decode response from server"
		}

		return se
	}

	return nil
}
