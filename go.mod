module github.com/eden-framework/eden-framework

go 1.14

replace k8s.io/client-go => k8s.io/client-go v0.18.8

require (
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/fatih/color v1.7.0
	github.com/go-courier/enumeration v1.0.1
	github.com/go-courier/metax v1.2.1
	github.com/go-courier/oas v1.2.0
	github.com/go-courier/ptr v1.0.1
	github.com/go-courier/reflectx v1.3.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.0
	github.com/imdario/mergo v0.3.9
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/gomega v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/profzone/envconfig v1.4.4
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/stretchr/testify v1.4.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	golang.org/x/tools v0.0.0-20200330040139-fa3cc9eebcfe
	google.golang.org/grpc v1.27.0
	gopkg.in/robfig/cron.v2 v2.0.0-20150107220207-be2e0b0deed5
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v12.0.0+incompatible
)
