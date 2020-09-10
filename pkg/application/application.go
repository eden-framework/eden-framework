package application

import (
	"github.com/profzone/eden-framework/internal"
	"github.com/profzone/eden-framework/internal/generator"
	"github.com/profzone/eden-framework/internal/project"
	"github.com/profzone/eden-framework/pkg/conf"
	"github.com/profzone/eden-framework/pkg/context"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/profzone/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
)

type Application struct {
	ctx                *context.WaitStopContext
	p                  *project.Project
	cmd                *cobra.Command
	envConfigPrefix    string
	outputDockerConfig bool
	autoMigration      bool
	Config             interface{}
}

func NewApplication(runner func(app *Application) error, config interface{}) *Application {
	p := &project.Project{}
	ctx := context.NewWaitStopContext()
	err := p.UnmarshalFromFile("", "")
	if err != nil {
		logrus.Panic(err)
	}

	tpe := reflect.TypeOf(config)
	if tpe.Kind() != reflect.Ptr {
		logrus.Panic("config must be a ptr value")
	}

	app := &Application{
		p:      p,
		ctx:    ctx,
		Config: config,
	}

	app.cmd = &cobra.Command{
		Use:   app.p.Name,
		Short: app.p.Desc,
		Run: func(cmd *cobra.Command, args []string) {
			go runner(app)
			app.WaitStop(func(ctx *context.WaitStopContext) error {
				ctx.Cancel()
				return nil
			})
		},
	}

	app.cmd.PersistentFlags().StringVarP(&app.envConfigPrefix, "env-prefix", "e", app.p.Name, "prefix for env var")
	app.cmd.PersistentFlags().BoolVarP(&app.outputDockerConfig, "docker", "d", true, "whether or not output configuration of docker")
	app.cmd.PersistentFlags().BoolVarP(&app.autoMigration, "db-migration", "m", os.Getenv("GOENV") == "DEV" || os.Getenv("GOENV") == "TEST", "auto migrate database if needed")

	return app
}

func (app *Application) AddCommand(cmd ...*cobra.Command) {
	app.cmd.AddCommand(cmd...)
}

func (app *Application) Start() {
	os.Setenv(internal.EnvVarKeyProjectName, app.p.Name)
	os.Setenv(internal.EnvVarKeyServiceName, strings.Replace(app.p.Name, "service-", "", 1))
	os.Setenv(internal.EnvVarKeyProjectGroup, app.p.Group)

	app.envConfigPrefix = str.ToUpperSnakeCase(app.envConfigPrefix)
	err := envconfig.Process(app.envConfigPrefix, app.Config)
	if err != nil {
		logrus.Panic(err)
	}
	envconfig.Usage(app.envConfigPrefix, app.Config)
	envVars, err := envconfig.GatherInfo(app.envConfigPrefix, app.Config)
	if err != nil {
		logrus.Panic(err)
	}

	if app.outputDockerConfig {
		generate := generator.NewDockerGenerator(app.p.Name, envVars)
		generator.Generate(generate, "", "")
	}

	// initialize global object
	conf.Initialize(app.Config)

	if err := app.cmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func (app *Application) WaitStop(clearFunc func(ctx *context.WaitStopContext) error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-sig
	err := clearFunc(app.ctx)
	if err != nil {
		logrus.Errorf("shutdown with error: %v", err)
	} else {
		logrus.Infof("graceful shutdown with signal: %s", s.String())
	}
}

func (app *Application) Context() *context.WaitStopContext {
	return app.ctx
}
