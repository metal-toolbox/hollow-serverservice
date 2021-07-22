package hollowserver

import (
	"net/http"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"

	"go.metalkube.net/hollow/internal/db"
	v1api "go.metalkube.net/hollow/pkg/api/v1"
	"go.metalkube.net/hollow/pkg/ginjwt"
)

// Server implements the Hollow server
type Server struct {
	Logger     *zap.Logger
	Listen     string
	Debug      bool
	Store      *db.Store
	AuthConfig AuthConfig
}

// AuthConfig provides the configuration for the authentication service
type AuthConfig struct {
	Audience string
	Issuer   string
	JWKSURI  string
}

var (
	readTimeout  = 10 * time.Second
	writeTimeout = 20 * time.Second
)

func (s *Server) setup() *gin.Engine {
	var (
		authMW *ginjwt.Middleware
		err    error
	)

	authMW, err = ginjwt.NewAuthMiddleware(s.AuthConfig.Audience, s.AuthConfig.Issuer, s.AuthConfig.JWKSURI)
	if err != nil {
		s.Logger.Sugar().Fatal("failed to initialize auth middleware", "error", err)
	}

	// Setup default gin router
	r := gin.New()

	// Health endpoints
	// These are defined before zap so that they are not logged
	r.GET("/healthz", s.livenessCheck)
	r.GET("/healthz/liveness", s.livenessCheck)
	r.GET("/healthz/readiness", s.readinessCheck)

	p := ginprometheus.NewPrometheus("gin")

	v1Rtr := v1api.Router{Store: s.Store, AuthMW: authMW}

	// Remove any params from the URL string to keep the number of labels down
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path

		for _, p := range c.Params {
			if p.Key == "uuid" {
				url = strings.Replace(url, p.Value, ":uuid", 1)
				break
			}
		}

		return url
	}

	p.Use(r)

	r.Use(ginzap.Logger(s.Logger, ginzap.WithTimeFormat(time.RFC3339),
		ginzap.WithUTC(true),
		ginzap.WithCustomFields(
			func(c *gin.Context) zap.Field { return zap.String("jwt_subject", ginjwt.GetSubject(c)) },
		),
	))
	r.Use(ginzap.RecoveryWithZap(s.Logger, true))

	v1 := r.Group("/api/v1")
	{
		v1Rtr.Routes(v1)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid request - route not found"})
	})

	return r
}

// NewServer returns a configured server
func (s *Server) NewServer() *http.Server {
	if !s.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return &http.Server{
		Handler:      s.setup(),
		Addr:         s.Listen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

// Run will start the server listening on the specified address
func (s *Server) Run() error {
	if !s.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return s.setup().Run(s.Listen)
}

// livenessCheck ensures that the server is up and responding
func (s *Server) livenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// readinessCheck ensures that the server is up and that we are able to process
// requests. Currently our only dependency is the DB so we just ensure that is
// responding.
func (s *Server) readinessCheck(c *gin.Context) {
	if s.Store.Ping() {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
		})
	}
}
