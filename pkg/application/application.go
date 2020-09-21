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
	outputK8sConfig    bool
	Config             []interface{}
}

func NewApplication(runner func(app *Application) error, config ...interface{}) *Application {
	p := &project.Project{}
	ctx := context.NewWaitStopContext()
	err := p.UnmarshalFromFile("", "")
	if err != nil {
		logrus.Panic(err)
	}

	for i, c := range config {
		tpe := reflect.TypeOf(c)
		if tpe.Kind() != reflect.Ptr {
			logrus.Panicf("the [%d] config must be a ptr value", i)
		}
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
	app.cmd.PersistentFlags().BoolVarP(&app.outputDockerConfig, "with-docker", "d", true, "whether or not output configuration of docker")
	app.cmd.PersistentFlags().BoolVarP(&app.outputK8sConfig, "with-k8s", "k", true, "whether or not output configuration of k8s")

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
	var envVars = make([]envconfig.EnvVar, 0)
	for _, c := range app.Config {
		err := envconfig.Process(app.envConfigPrefix, c)
		if err != nil {
			logrus.Panic(err)
		}
		envconfig.Usage(app.envConfigPrefix, c)

		envs, err := envconfig.GatherInfo(app.envConfigPrefix, c)
		if err != nil {
			logrus.Panic(err)
		}

		envVars = append(envVars, envs...)
	}

	if os.Getenv("GOENV") != "PROD" {
		cwd, _ := os.Getwd()

		generate := generator.NewDockerGenerator(app.p.Name, envVars)
		generator.Generate(generate, cwd, cwd)

		k8sGenerator := generator.NewK8sGenerator(app.Config)
		generator.Generate(k8sGenerator, cwd, cwd)
	}

	// initialize global object
	conf.Initialize(app.Config...)

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
