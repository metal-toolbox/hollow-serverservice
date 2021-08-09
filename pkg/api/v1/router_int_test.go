package hollow_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"go.metalkube.net/hollow/internal/db"
	"go.metalkube.net/hollow/internal/hollowserver"
	hollow "go.metalkube.net/hollow/pkg/api/v1"
	"go.metalkube.net/hollow/pkg/ginjwt"
)

type integrationServer struct {
	h      http.Handler
	Client *hollow.Client
}

func serverTest(t *testing.T) *integrationServer {
	jwksURI := ginjwt.TestHelperJWKSProvider()

	store := db.DatabaseTest(t)

	l, _ := zap.NewDevelopment()

	hs := hollowserver.Server{
		Logger: l,
		Store:  store,
		AuthConfig: ginjwt.AuthConfig{
			Enabled:  true,
			Audience: "hollow.test",
			Issuer:   "hollow.test.issuer",
			JWKSURI:  jwksURI,
		},
	}
	s := hs.NewServer()

	ts := &integrationServer{
		h: s.Handler,
	}

	c, err := hollow.NewClient("testToken", "http://test.hollow.com", ts)
	require.NoError(t, err)

	ts.Client = c

	return ts
}

func (s *integrationServer) Do(req *http.Request) (*http.Response, error) {
	// if the context is expired return the error, used for timeout tests
	if err := req.Context().Err(); err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	s.h.ServeHTTP(w, req)

	return w.Result(), nil
}

func validToken(scopes []string) string {
	claims := jwt.Claims{
		Subject:   "test-user",
		Issuer:    "hollow.test.issuer",
		NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		Expiry:    jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  jwt.Audience{"hollow.test", "another.test.service"},
	}
	signer := ginjwt.TestHelperMustMakeSigner(jose.RS256, ginjwt.TestPrivRSAKey1ID, ginjwt.TestPrivRSAKey1)

	return ginjwt.TestHelperGetToken(signer, claims, scopes)
}
