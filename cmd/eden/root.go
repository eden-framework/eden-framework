package main

import (
	"github.com/eden-framework/eden-framework/internal/project"
	_ "github.com/eden-framework/eden-framework/internal/workflows"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var currentProject = &project.Project{}

func init() {
	_ = currentProject.UnmarshalFromFile("", "")

	project.RegisterBuilder("BUILDER_RANCHER", &project.Builder{
		Image:      getEnvOrDefault("BUILDER_RANCHER", "profzone/golang-ondeploy:2.4.3"),
		WorkingDir: "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.RegisterBuilder("BUILDER_DOCKER", &project.Builder{
		Image:      getEnvOrDefault("BUILDER_DOCKER", "profzone/golang-onship:1.14"),
		WorkingDir: "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.RegisterBuilder("BUILDER_GOLANG", &project.Builder{
		ProgramLanguage: "golang",
		Image:           getEnvOrDefault("BUILDER_GOLANG", "profzone/golang-onbuild:1.14"),
		WorkingDir:      "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})

	if currentProject.Scripts != nil {
		for scriptCmd, script := range currentProject.Scripts {
			ciRunCmd.AddCommand(&cobra.Command{
				Use:   scriptCmd,
				Short: script.String(),
				Run: func(cmd *cobra.Command, args []string) {
					err := currentProject.RunScript(scriptCmd, ciRunCmdInDocker)
					if err != nil {
						logrus.Error(err)
					}
				},
			})
		}
	}
}

func getEnvOrDefault(key string, value string) string {
	envVar := os.Getenv(key)
	if envVar != "" {
		return envVar
	}
	return value
}
