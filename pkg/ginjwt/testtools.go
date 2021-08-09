//+build testtools

package ginjwt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	testKeySize = 2048

	// TestPrivRSAKey1 provides an RSA key used to sign tokens
	TestPrivRSAKey1, _ = rsa.GenerateKey(rand.Reader, testKeySize)
	// TestPrivRSAKey1ID is the ID of this signing key in tokens
	TestPrivRSAKey1ID = "testKey1"
	// TestPrivRSAKey2 provides an RSA key used to sign tokens
	TestPrivRSAKey2, _ = rsa.GenerateKey(rand.Reader, testKeySize)
	// TestPrivRSAKey2ID is the ID of this signing key in tokens
	TestPrivRSAKey2ID = "testKey2"
)

// TestHelperMustMakeSigner will return a JWT signer from the given key
func TestHelperMustMakeSigner(alg jose.SignatureAlgorithm, kid string, k interface{}) jose.Signer {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: alg, Key: k}, (&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", kid))
	if err != nil {
		panic("failed to create signer:" + err.Error())
	}

	return sig
}

// TestHelperJWKSProvider returns a url for a webserver that will return JSONWebKeySets
func TestHelperJWKSProvider() string {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/.well-known/jwks.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				{
					KeyID: TestPrivRSAKey1ID,
					Key:   &TestPrivRSAKey1.PublicKey,
				},
				{
					KeyID: TestPrivRSAKey2ID,
					Key:   &TestPrivRSAKey2.PublicKey,
				},
			},
		})
	})

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	s := &http.Server{
		Handler: r,
	}

	go func() {
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", listener.Addr().(*net.TCPAddr).Port)
}

// TestHelperGetToken will return a signed token
func TestHelperGetToken(signer jose.Signer, cl jwt.Claims, scopes []string) string {
	sc := map[string]interface{}{}

	sc["scope"] = strings.Join(scopes, " ")

	raw, err := jwt.Signed(signer).Claims(cl).Claims(sc).CompactSerialize()
	if err != nil {
		panic(err)
	}

	return raw
}
