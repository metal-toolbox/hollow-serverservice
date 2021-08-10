package ginjwt_test

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"go.metalkube.net/hollow/pkg/ginjwt"
)

func TestMiddlewareValidatesTokens(t *testing.T) {
	var testCases = []struct {
		testName         string
		middlewareAud    string
		middlewareIss    string
		middlewareScopes []string
		signingKey       *rsa.PrivateKey
		signingKeyID     string
		claims           jwt.Claims
		claimScopes      []string
		responseCode     int
		responseBody     string
	}{
		{
			"unknown keyid",
			"ginjwt.test",
			"ginjwt.test.issuer2",
			[]string{"testScope"},
			ginjwt.TestPrivRSAKey1,
			"randomUnknownID",
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusUnauthorized,
			"invalid token signing key",
		},
		{
			"incorrect keyid",
			"ginjwt.test",
			"ginjwt.test.issuer2",
			[]string{"testScope"},
			ginjwt.TestPrivRSAKey1,
			ginjwt.TestPrivRSAKey2ID,
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusUnauthorized,
			"unable to validate auth token",
		},
		{
			"incorrect issuer",
			"ginjwt.test",
			"ginjwt.test.issuer2",
			[]string{"testScope"},
			ginjwt.TestPrivRSAKey1,
			ginjwt.TestPrivRSAKey1ID,
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusUnauthorized,
			"invalid issuer claim",
		},
		{
			"incorrect audience",
			"ginjwt.testFail",
			"ginjwt.test.issuer",
			[]string{"testScope"},
			ginjwt.TestPrivRSAKey1,
			ginjwt.TestPrivRSAKey1ID,
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusUnauthorized,
			"invalid audience claim",
		},
		{
			"incorrect scopes",
			"ginjwt.test",
			"ginjwt.test.issuer",
			[]string{"adminscope"},
			ginjwt.TestPrivRSAKey1,
			ginjwt.TestPrivRSAKey1ID,
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusForbidden,
			"missing required scope",
		},
		{
			"expired token",
			"ginjwt.test",
			"ginjwt.test.issuer",
			[]string{"testScope"},
			ginjwt.TestPrivRSAKey1,
			ginjwt.TestPrivRSAKey1ID,
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(-6 * time.Hour)),
				Expiry:    jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusUnauthorized,
			"token is expired",
		},
		{
			"future token",
			"ginjwt.test",
			"ginjwt.test.issuer",
			[]string{"testScope"},
			ginjwt.TestPrivRSAKey1,
			ginjwt.TestPrivRSAKey1ID,
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusUnauthorized,
			"token not valid yet",
		},
		{
			"happy path",
			"ginjwt.test",
			"ginjwt.test.issuer",
			[]string{"testScope"},
			ginjwt.TestPrivRSAKey1,
			ginjwt.TestPrivRSAKey1ID,
			jwt.Claims{
				Subject:   "test-user",
				Issuer:    "ginjwt.test.issuer",
				NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				Audience:  jwt.Audience{"ginjwt.test", "another.test.service"},
			},
			[]string{"testScope", "anotherScope", "more-scopes"},
			http.StatusOK,
			"ok",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			jwksURI := ginjwt.TestHelperJWKSProvider()

			cfg := ginjwt.AuthConfig{Enabled: true, Audience: tt.middlewareAud, Issuer: tt.middlewareIss, JWKSURI: jwksURI}
			authMW, err := ginjwt.NewAuthMiddleware(cfg)
			require.NoError(t, err)

			r := gin.New()
			r.Use(authMW.AuthRequired(tt.middlewareScopes))
			r.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, "ok")
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "http://test/", nil)

			signer := ginjwt.TestHelperMustMakeSigner(jose.RS256, tt.signingKeyID, tt.signingKey)
			rawToken := ginjwt.TestHelperGetToken(signer, tt.claims, "scope", strings.Join(tt.claimScopes, " "))
			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", rawToken))

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.responseCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestInvalidAuthHeader(t *testing.T) {
	var testCases = []struct {
		testName         string
		authHeader       string
		responseCode     int
		responseContains string
	}{
		{
			"no auth header",
			"",
			http.StatusUnauthorized,
			"missing authorization header",
		},
		{
			"wrong format",
			"notbearer token",
			http.StatusUnauthorized,
			"invalid authorization header",
		},
		{
			"invalid token",
			"bearer token",
			http.StatusUnauthorized,
			"unable to parse auth token",
		},
		{
			"token with no kid",
			"bearer eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJpc3N1ZXIiLCJzY29wZXMiOlsiczEiLCJzMiJdLCJzdWIiOiJzdWJqZWN0In0.UDDtyK9gC9kyHltcP7E_XODsnqcJWZIiXeGmSAH7SE9YKy3N0KSfFIN85dCNjTfs6zvy4rkrCHzLB7uKAtzMearh3q7jL4nxbhUMhlUcs_9QDVoN4q_j58XmRqBqRnBk-RmDu9TgcV8RbErP4awpIhwWb5UU-hR__4_iNbHdKqwSUPDKYGlf5eicuiYrPxH8mxivk4LRD-vyRdBZZKBt0XIDnEU4TdcNCzAXojkftqcFWYsczwS8R4JHd1qYsMyiaWl4trdHZkO4QkeLe34z4ZAaPMt3wE-gcU-VoqYTGxz-K3Le2VaZ0r3j_z6bOInsv0yngC_cD1dCXMyQJWnWjQ",
			http.StatusUnauthorized,
			"unable to parse auth token header",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			jwksURI := ginjwt.TestHelperJWKSProvider()
			cfg := ginjwt.AuthConfig{Enabled: true, Audience: "aud", Issuer: "iss", JWKSURI: jwksURI}
			authMW, err := ginjwt.NewAuthMiddleware(cfg)
			require.NoError(t, err)

			r := gin.New()
			r.Use(authMW.AuthRequired([]string{"auth"}))
			r.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, "ok")
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "http://test/", nil)

			req.Header.Set("Authorization", tt.authHeader)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.responseCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.responseContains)
		})
	}
}
