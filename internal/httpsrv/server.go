package httpsrv

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.hollow.sh/toolbox/ginjwt"
	"go.hollow.sh/toolbox/version"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gocloud.dev/secrets"

	v1api "go.hollow.sh/serverservice/pkg/api/v1"
)

// Server implements the HTTP Server
type Server struct {
	Logger        *zap.Logger
	Listen        string
	Debug         bool
	DB            *sqlx.DB
	AuthConfig    ginjwt.AuthConfig
	SecretsKeeper *secrets.Keeper
}

var (
	readTimeout  = 10 * time.Second
	writeTimeout = 20 * time.Second
	corsMaxAge   = 12 * time.Hour
)

func (s *Server) setup() *gin.Engine {
	var (
		authMW *ginjwt.Middleware
		err    error
	)

	authMW, err = ginjwt.NewAuthMiddleware(s.AuthConfig)
	if err != nil {
		s.Logger.Sugar().Fatal("failed to initialize auth middleware", "error", err)
	}

	// Setup default gin router
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           corsMaxAge,
	}))

	p := ginprometheus.NewPrometheus("gin")

	v1Rtr := v1api.Router{
		DB:            s.DB,
		AuthMW:        authMW,
		SecretsKeeper: s.SecretsKeeper,
	}

	// Remove any params from the URL string to keep the number of labels down
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}

	p.Use(r)

	r.Use(ginzap.Logger(s.Logger.With(zap.String("component", "httpsrv")), ginzap.WithTimeFormat(time.RFC3339),
		ginzap.WithUTC(true),
		ginzap.WithCustomFields(
			func(c *gin.Context) zap.Field { return zap.String("jwt_subject", ginjwt.GetSubject(c)) },
			func(c *gin.Context) zap.Field { return zap.String("jwt_user", ginjwt.GetUser(c)) },
		),
	))
	r.Use(ginzap.RecoveryWithZap(s.Logger.With(zap.String("component", "httpsrv")), true))

	tp := otel.GetTracerProvider()
	if tp != nil {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		r.Use(otelgin.Middleware(hostname, otelgin.WithTracerProvider(tp)))
	}

	// Version endpoint returns build information
	r.GET("/version", s.version)

	// Health endpoints
	r.GET("/healthz", s.livenessCheck)
	r.GET("/healthz/liveness", s.livenessCheck)
	r.GET("/healthz/readiness", s.readinessCheck)

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
	if err := s.DB.PingContext(c.Request.Context()); err != nil {
		s.Logger.Sugar().Errorf("readiness check db ping failed", "err", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// version returns the serverservice build information.
func (s *Server) version(c *gin.Context) {
	c.JSON(http.StatusOK, version.String())
}
