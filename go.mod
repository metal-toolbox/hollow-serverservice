module go.hollow.sh/serverservice

go 1.17

require (
	github.com/XSAM/otelsql v0.5.0
	github.com/cockroachdb/cockroach-go/v2 v2.1.1
	github.com/friendsofgo/errors v0.9.2
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.7.4
	github.com/google/uuid v1.2.0
	github.com/gosimple/slug v1.10.0
	github.com/lib/pq v1.10.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.1.0
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/sqlboiler/v4 v4.6.0
	github.com/volatiletech/strmangle v0.0.1
	github.com/zsais/go-gin-prometheus v0.1.0
	go.hollow.sh/toolbox v0.0.0-20210826144247-5ed6c7643625
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.uber.org/zap v1.17.0
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
)

require (
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	gopkg.in/ini.v1 v1.63.1 // indirect
)

replace github.com/gin-contrib/zap => github.com/thinkgos/zap v0.0.2-0.20210226022008-5b2cf0c4d297
