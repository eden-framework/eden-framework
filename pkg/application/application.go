package application

import (
	"github.com/profzone/eden-framework/internal"
	"github.com/profzone/eden-framework/internal/project"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type Application struct {
	p                  *project.Project
	envConfigPrefix    string
	outputDockerConfig bool
	autoMigration      bool
	Runner             func() error
}

func NewApplication(runner func() error) *Application {
	p := &project.Project{}
	err := p.UnmarshalFromFile("", "")
	if err != nil {
		logrus.Panic(err)
	}

	return &Application{
		p:      p,
		Runner: runner,
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

	command.PersistentFlags().StringVarP(&app.envConfigPrefix, "env-prefix", "e", "S", "prefix for env var")
	command.PersistentFlags().BoolVarP(&app.outputDockerConfig, "docker", "d", true, "whether or not output configuration of docker")
	command.PersistentFlags().BoolVarP(&app.autoMigration, "db-migration", "m", os.Getenv("GOENV") == "DEV" || os.Getenv("GOENV") == "TEST", "auto migrate database if needed")

	if err := command.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	if err := app.Runner(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
