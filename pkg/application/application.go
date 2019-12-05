package application

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/profzone/eden-framework/internal"
	"github.com/profzone/eden-framework/internal/project"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"strings"
)

type Application struct {
	p                  *project.Project
	envConfigPrefix    string
	outputDockerConfig bool
	autoMigration      bool
	Runner             func(app *Application) error
	Config             interface{}
}

func NewApplication(runner func(app *Application) error, config interface{}) *Application {
	p := &project.Project{}
	err := p.UnmarshalFromFile("", "")
	if err != nil {
		logrus.Panic(err)
	}

	tpe := reflect.TypeOf(config)
	if tpe.Kind() != reflect.Ptr {
		logrus.Panic("config must be a ptr value")
	}

	return &Application{
		p:      p,
		Runner: runner,
		Config: config,
	}
}

func (app *Application) Start() {
	os.Setenv(internal.EnvVarKeyProjectName, app.p.Name)
	os.Setenv(internal.EnvVarKeyServiceName, strings.Replace(app.p.Name, "service-", "", 1))
	os.Setenv(internal.EnvVarKeyProjectGroup, app.p.Group)

	command := &cobra.Command{
		Use:   app.p.Name,
		Short: app.p.Desc,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	command.PersistentFlags().StringVarP(&app.envConfigPrefix, "env-prefix", "e", app.p.Name, "prefix for env var")
	command.PersistentFlags().BoolVarP(&app.outputDockerConfig, "docker", "d", true, "whether or not output configuration of docker")
	command.PersistentFlags().BoolVarP(&app.autoMigration, "db-migration", "m", os.Getenv("GOENV") == "DEV" || os.Getenv("GOENV") == "TEST", "auto migrate database if needed")
	app.envConfigPrefix = str.ToUpperSnakeCase(app.envConfigPrefix)

	if err := command.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	err := envconfig.Usage(app.envConfigPrefix, app.Config)
	if err != nil {
		logrus.Panic(err)
	}
	err = envconfig.Process(app.envConfigPrefix, app.Config)
	if err != nil {
		logrus.Panic(err)
	}

	if err := app.Runner(app); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
