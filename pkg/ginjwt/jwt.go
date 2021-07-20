package ginjwt

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const contextKeySubject = "jwt.subject"

// Middleware provides a gin compatible middleware that will authenticate JWT requests
type Middleware struct {
	audience   string
	issuer     string
	jwksURI    string
	cachedJWKS jose.JSONWebKeySet
}

type customClaims struct {
	Scope string `json:"scope"`
}

func (c *customClaims) Scopes() []string {
	return strings.Split(c.Scope, " ")
}

// NewAuthMiddleware will return an auth middleware configured with the jwt parameters passed in
func NewAuthMiddleware(aud, iss, jwksURI string) (*Middleware, error) {
	mw := &Middleware{
		audience: aud,
		issuer:   iss,
		jwksURI:  jwksURI,
	}

	if err := mw.refreshJWKS(); err != nil {
		return nil, err
	}

	return mw, nil
}

// AuthRequired provides a middleware that ensures a request has authentication
func (m *Middleware) AuthRequired(scopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeaderParts := strings.Split(c.Request.Header.Get("Authorization"), " ")
		rawToken := authHeaderParts[1]

		tok, err := jwt.ParseSigned(rawToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unable to parse auth token"})
			return
		}

		if tok.Headers[0].KeyID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unable to parse auth token header"})
			return
		}

		key := m.getJWKS(tok.Headers[0].KeyID)
		if key == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token signing key"})
			return
		}

		cl := jwt.Claims{}
		sc := customClaims{}

		if err := tok.Claims(key, &cl, &sc); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unable to validate auth token"})
			return
		}

		err = cl.Validate(jwt.Expected{
			Issuer:   m.issuer,
			Audience: jwt.Audience{m.audience},
			Time:     time.Now(),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid auth token", "error": err.Error()})
			return
		}

		if !hasScope(sc.Scopes(), scopes) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized, missing required scope"})
			return
		}

		c.Set(contextKeySubject, cl.Subject)
	}
}

func (m *Middleware) refreshJWKS() error {
	resp, err := http.Get(m.jwksURI) //nolint:noctx
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&m.cachedJWKS)
}

func (m *Middleware) getJWKS(kid string) *jose.JSONWebKey {
	keys := m.cachedJWKS.Key(kid)
	if len(keys) == 0 {
		// couldn't find the signing key in our cache, refresh cache and search again
		if err := m.refreshJWKS(); err != nil {
			return nil
		}

		keys = m.cachedJWKS.Key(kid)
		if len(keys) == 0 {
			return nil
		}
	}

	return &keys[0]
}

func hasScope(have, needed []string) bool {
	neededMap := make(map[string]bool)
	for _, s := range needed {
		neededMap[s] = true
	}

	for _, s := range have {
		if neededMap[s] {
			return true
		}
	}

	return false
}

// GetSubject will return the JWT subject that is saved in the request. This requires that authentication of the request
// has already occurred. If authentication failed or there isn't a user, an empty string is returned. This returns
// whatever value was in the JWT subject field and might not be a human readable value
func GetSubject(c *gin.Context) string {
	return c.GetString(contextKeySubject)
}
