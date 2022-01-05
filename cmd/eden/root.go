package main

import (
	"gitee.com/eden-framework/eden-framework/internal/project"
	_ "gitee.com/eden-framework/eden-framework/internal/workflows"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var currentProject = &project.Project{}

func init() {
	_ = currentProject.UnmarshalFromFile("", "")

	if currentProject.Scripts != nil {
		for scriptCmd, script := range currentProject.Scripts {
			ciRunCmd.AddCommand(&cobra.Command{
				Use:   scriptCmd,
				Short: script.String(),
				Run: func(cmd *cobra.Command, args []string) {
					err := currentProject.RunScript(cmd.Use, ciRunCmdInDocker)
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

func getConfigOrDefault(key string, value string) string {
	val := viper.GetString(key)
	if val != "" {
		return val
	}
	return value
}
