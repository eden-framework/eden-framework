package main

import (
	"github.com/profzone/eden-framework/internal/project"
	_ "github.com/profzone/eden-framework/internal/workflows"
	"os"
)

var currentProject = &project.Project{}

func init() {
	_ = currentProject.UnmarshalFromFile("", "")

	project.RegisterBuilder("BUILDER_RANCHER", &project.Builder{
		Image:      getEnvOrDefault("BUILDER_RANCHER", "g7pay/env-rancher-cli:latest"),
		WorkingDir: "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.RegisterBuilder("BUILDER_GOLANG", &project.Builder{
		ProgramLanguage: "golang",
		Image:           getEnvOrDefault("BUILDER_GOLANG", "profzone/golang:onbuild"),
		WorkingDir:      "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.RegisterBuilder("BUILDER_VUE", &project.Builder{
		ProgramLanguage: "vue",
		Image:           getEnvOrDefault("BUILDER_VUE", "profzone/node:alpine"),
		WorkingDir:      "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.RegisterBuilder("BUILDER_GITBOOK", &project.Builder{
		ProgramLanguage: "gitbook",
		Image:           getEnvOrDefault("BUILDER_GITBOOK", "g7pay/env-node:gitbook-builder"),
		WorkingDir:      "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
}

func getEnvOrDefault(key string, value string) string {
	envVar := os.Getenv(key)
	if envVar != "" {
		return envVar
	}
	return value
}
