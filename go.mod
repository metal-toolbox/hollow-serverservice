module go.metalkube.net/hollow

go 1.16

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.7.1
	github.com/google/uuid v1.2.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gosimple/slug v1.10.0
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.1.13 // indirect
	github.com/zsais/go-gin-prometheus v0.1.0
	go.uber.org/zap v1.16.0
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.9
	gorm.io/plugin/prometheus v0.0.0-20210614014227-3996fd54c851
	moul.io/zapgorm2 v1.0.1
)

replace github.com/gin-contrib/zap => github.com/thinkgos/zap v0.0.2-0.20210226022008-5b2cf0c4d297
