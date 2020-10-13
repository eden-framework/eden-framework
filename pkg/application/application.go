package application

import (
	"github.com/eden-framework/eden-framework/internal"
	"github.com/eden-framework/eden-framework/internal/generator"
	"github.com/eden-framework/eden-framework/internal/project"
	"github.com/eden-framework/eden-framework/pkg/conf"
	"github.com/eden-framework/eden-framework/pkg/conf/apollo"
	"github.com/eden-framework/eden-framework/pkg/context"
	str "github.com/eden-framework/strings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Application struct {
	ctx             *context.WaitStopContext
	p               *project.Project
	cmd             *cobra.Command
	envConfigPrefix string
	apolloConfig    *apollo.ApolloBaseConfig
	envConfig       []interface{}
}

func NewApplication(runner func(ctx *context.WaitStopContext) error, opts ...Option) *Application {
	p := &project.Project{}
	ctx := context.NewWaitStopContext()
	err := p.UnmarshalFromFile("", "")
	if err != nil {
		logrus.Panic(err)
	}

	app := &Application{
		p:   p,
		ctx: ctx,
	}

	for _, opt := range opts {
		opt(app)
	}

	app.cmd = &cobra.Command{
		Use:   app.p.Name,
		Short: app.p.Desc,
		Run: func(cmd *cobra.Command, args []string) {
			go runner(ctx)
			app.WaitStop(func(ctx *context.WaitStopContext) error {
				ctx.Cancel()
				return nil
			})
		},
	}

	app.cmd.PersistentFlags().StringVarP(&app.envConfigPrefix, "env-prefix", "e", app.p.Name, "prefix for env var")

	return app
}

func (app *Application) AddCommand(cmd ...*cobra.Command) {
	app.cmd.AddCommand(cmd...)
}

func (app *Application) Start() {
	app.envConfigPrefix = str.ToUpperSnakeCase(app.envConfigPrefix)

	os.Setenv(internal.EnvVarKeyProjectName, app.p.Name)
	os.Setenv(internal.EnvVarKeyServiceName, strings.Replace(app.p.Name, "service-", "", 1))
	os.Setenv(internal.EnvVarKeyProjectGroup, app.p.Group)

	// config from env
	var confs []interface{}
	if app.apolloConfig != nil && os.Getenv("GOENV") != "LOCAL" {
		confs = append(confs, app.apolloConfig)
	} else {
		confs = append(confs, app.envConfig...)
	}
	envVars := conf.FromEnv(app.envConfigPrefix, confs)

	// config from apollo
	conf.FromApollo(app.apolloConfig, app.envConfig)

	// output config env variables
	if os.Getenv("GOENV") != "PROD" {
		cwd, _ := os.Getwd()

		generate := generator.NewDockerGenerator(app.p.Name, envVars)
		generator.Generate(generate, cwd, cwd)

		k8sGenerator := generator.NewK8sGenerator(app.envConfig)
		generator.Generate(k8sGenerator, cwd, cwd)
	}

	// initialize global object
	conf.Initialize(app.envConfig...)

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
