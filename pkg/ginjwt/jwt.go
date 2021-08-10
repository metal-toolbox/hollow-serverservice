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

const (
	contextKeySubject       = "jwt.subject"
	contextKeyUser          = "jwt.user"
	expectedAuthHeaderParts = 2
)

// Middleware provides a gin compatible middleware that will authenticate JWT requests
type Middleware struct {
	config     AuthConfig
	cachedJWKS jose.JSONWebKeySet
}

// AuthConfig provides the configuration for the authentication service
type AuthConfig struct {
	Enabled       bool
	Audience      string
	Issuer        string
	JWKSURI       string
	LogFields     []string
	RolesClaim    string
	UsernameClaim string
}

// NewAuthMiddleware will return an auth middleware configured with the jwt parameters passed in
func NewAuthMiddleware(cfg AuthConfig) (*Middleware, error) {
	if cfg.RolesClaim == "" {
		cfg.RolesClaim = "scope"
	}

	if cfg.UsernameClaim == "" {
		cfg.UsernameClaim = "sub"
	}

	mw := &Middleware{
		config: cfg,
	}

	if !cfg.Enabled {
		return mw, nil
	}

	if err := mw.refreshJWKS(); err != nil {
		return nil, err
	}

	return mw, nil
}

// AuthRequired provides a middleware that ensures a request has authentication
func (m *Middleware) AuthRequired(scopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			return
		}

		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header, expected format: \"Bearer token\""})
			return
		}

		authHeaderParts := strings.SplitN(authHeader, " ", expectedAuthHeaderParts)

		if !(len(authHeaderParts) == expectedAuthHeaderParts && strings.ToLower(authHeaderParts[0]) == "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header, expected format: \"Bearer token\""})
			return
		}

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
		sc := map[string]interface{}{}

		if err := tok.Claims(key, &cl, &sc); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unable to validate auth token"})
			return
		}

		err = cl.Validate(jwt.Expected{
			Issuer:   m.config.Issuer,
			Audience: jwt.Audience{m.config.Audience},
			Time:     time.Now(),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid auth token", "error": err.Error()})
			return
		}

		var roles []string
		switch r := sc[m.config.RolesClaim].(type) {
		case string:
			roles = strings.Split(r, " ")
		case []interface{}:
			for _, i := range r {
				roles = append(roles, i.(string))
			}
		}

		if !hasScope(roles, scopes) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "not authorized, missing required scope"})
			return
		}

		var user string
		switch u := sc[m.config.UsernameClaim].(type) {
		case string:
			user = u
		default:
			user = cl.Subject
		}

		c.Set(contextKeySubject, cl.Subject)
		c.Set(contextKeyUser, user)
	}
}

func (m *Middleware) refreshJWKS() error {
	resp, err := http.Get(m.config.JWKSURI) //nolint:noctx
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

// GetUser will return the JWT user that is saved in the request. This requires that authentication of the request
// has already occurred. If authentication failed or there isn't a user an empty string is returned.
func GetUser(c *gin.Context) string {
	return c.GetString(contextKeyUser)
}
