package hollowserver

import (
	"net/http"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"

	v1api "go.metalkube.net/hollow/pkg/api/v1"
)

// Server implements the Hollow server
type Server struct {
	Logger *zap.Logger
	Listen string
}

var (
	readTimeout  = 10 * time.Second
	writeTimeout = 20 * time.Second
)

func (s *Server) setup() http.Handler {
	// Setup default gin router
	r := gin.New()

	p := ginprometheus.NewPrometheus("gin")

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

	r.Use(ginzap.Ginzap(s.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(s.Logger, true))

	v1 := r.Group("/api/v1")
	{
		v1api.RouteMap(v1)
	}

	// Add endpoints
	// router.GET("/status", status)

	return r
}

// NewServer returns a configured server
func (s *Server) NewServer() *http.Server {
	return &http.Server{
		Handler:      s.setup(),
		Addr:         s.Listen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}
