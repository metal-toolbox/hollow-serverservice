package dcimserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"go.hollow.sh/dcim/internal/dbtools"
	"go.hollow.sh/dcim/internal/dcimserver"
	"go.hollow.sh/dcim/pkg/ginjwt"
)

var serverAuthConfig = ginjwt.AuthConfig{
	Enabled: false,
}

func TestUnknownRoute(t *testing.T) {
	hs := dcimserver.Server{Logger: zap.NewNop(), AuthConfig: serverAuthConfig}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/a/route/that/doesnt/exist", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"message":"invalid request - route not found"}`, w.Body.String())
}

func TestHealthzRoute(t *testing.T) {
	hs := dcimserver.Server{Logger: zap.NewNop(), AuthConfig: serverAuthConfig}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"UP"}`, w.Body.String())
}

func TestLivenessRoute(t *testing.T) {
	hs := dcimserver.Server{Logger: zap.NewNop(), AuthConfig: serverAuthConfig}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/healthz/liveness", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"UP"}`, w.Body.String())
}

func TestReadinessRouteDown(t *testing.T) {
	hs := dcimserver.Server{Logger: zap.NewNop(), AuthConfig: serverAuthConfig}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/healthz/readiness", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 503, w.Code)
	assert.Equal(t, `{"status":"DOWN"}`, w.Body.String())
}

func TestReadinessRouteUp(t *testing.T) {
	db := dbtools.DatabaseTest(t)

	hs := dcimserver.Server{Logger: zap.NewNop(), AuthConfig: serverAuthConfig, DB: db}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/healthz/readiness", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"UP"}`, w.Body.String())
}
