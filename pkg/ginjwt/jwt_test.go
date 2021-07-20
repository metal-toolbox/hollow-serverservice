package ginjwt_test

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
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
			http.StatusUnauthorized,
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
			authMW, err := ginjwt.NewAuthMiddleware(tt.middlewareAud, tt.middlewareIss, jwksURI)
			require.NoError(t, err)

			r := gin.New()
			r.Use(authMW.AuthRequired(tt.middlewareScopes))
			r.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, "ok")
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "http://test/", nil)

			signer := ginjwt.TestHelperMustMakeSigner(jose.RS256, tt.signingKeyID, tt.signingKey)
			rawToken := ginjwt.TestHelperGetToken(signer, tt.claims, tt.claimScopes)
			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", rawToken))

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.responseCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.responseBody)
		})
	}
}
