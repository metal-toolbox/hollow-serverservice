module go.hollow.sh/serverservice

go 1.17

require (
	github.com/XSAM/otelsql v0.14.1
	github.com/cockroachdb/cockroach-go/v2 v2.2.8
	github.com/friendsofgo/errors v0.9.2
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/zap v0.0.2
	github.com/gin-gonic/gin v1.8.1
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.12.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.5.3
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.2
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/sqlboiler/v4 v4.11.0
	github.com/volatiletech/strmangle v0.0.4
	github.com/zsais/go-gin-prometheus v0.1.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.32.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.uber.org/zap v1.21.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
)

require go.hollow.sh/toolbox v0.1.4

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.8.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.0.6 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.6.2 // indirect
	github.com/jackc/pgx/v4 v4.10.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gin-contrib/zap => github.com/thinkgos/zap v0.0.2-0.20210226022008-5b2cf0c4d297
