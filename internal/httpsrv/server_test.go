package httpsrv_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.hollow.sh/toolbox/ginjwt"
	"go.uber.org/zap"

	"go.hollow.sh/serverservice/internal/dbtools"
	"go.hollow.sh/serverservice/internal/httpsrv"
)

var serverAuthConfig = []ginjwt.AuthConfig{
	{
		Enabled: false,
	},
}

func TestUnknownRoute(t *testing.T) {
	hs := httpsrv.Server{Logger: zap.NewNop(), AuthConfigs: serverAuthConfig}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/a/route/that/doesnt/exist", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"message":"invalid request - route not found"}`, w.Body.String())
}

func TestHealthzRoute(t *testing.T) {
	hs := httpsrv.Server{Logger: zap.NewNop(), AuthConfigs: serverAuthConfig}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"UP"}`, w.Body.String())
}

func TestLivenessRoute(t *testing.T) {
	hs := httpsrv.Server{Logger: zap.NewNop(), AuthConfigs: serverAuthConfig}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/healthz/liveness", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"UP"}`, w.Body.String())
}

func TestReadinessRouteDown(t *testing.T) {
	db, _ := sqlx.Open("postgres", "localhost:12341")

	hs := httpsrv.Server{Logger: zap.NewNop(), AuthConfigs: serverAuthConfig, DB: db}
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

	hs := httpsrv.Server{Logger: zap.NewNop(), AuthConfigs: serverAuthConfig, DB: db}
	s := hs.NewServer()
	router := s.Handler

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/healthz/readiness", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"UP"}`, w.Body.String())
}
