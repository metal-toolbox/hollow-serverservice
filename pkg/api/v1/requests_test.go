package fleetdbapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Do(t *testing.T) {
	serveMux1 := http.NewServeMux()

	serveMux1.HandleFunc(
		"/test",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				assert.Equal(t, "bearer dummy", r.Header.Get("Authorization"))

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"foo": "bar"}`))
			default:
				t.Fatal("expected GET request, got: " + r.Method)
			}
		},
	)

	serveMux2 := http.NewServeMux()

	serveMux2.HandleFunc(
		"/brokenjson",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"foo": "bar"}\n{"foo": "bar"}`))
			default:
				t.Fatal("expected GET request, got: " + r.Method)
			}
		},
	)

	testcases := []struct {
		name                  string
		endpoint              string
		expectedResult        interface{}
		serveMux              *http.ServeMux
		expectedErrorContains string
	}{
		{
			"happy path",
			"/test",
			map[string]string{"foo": "bar"},
			serveMux1,
			"",
		},
		{
			"broken json returns error",
			"/brokenjson",
			map[string]string{"foo": "bar"},
			serveMux2,
			"invalid character",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.serveMux)
			client, err := NewClientWithToken("dummy", server.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequestWithContext(context.TODO(), "GET", server.URL+tc.endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			result := map[string]string{}

			err = client.do(req, &result)
			if tc.expectedErrorContains != "" {
				if err == nil {
					t.Fatalf("expected error: '%s', got nil", tc.expectedErrorContains)
				}

				assert.Contains(t, err.Error(), tc.expectedErrorContains)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
