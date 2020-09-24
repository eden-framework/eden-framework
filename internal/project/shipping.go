package project

import (
	"github.com/eden-framework/eden-framework/internal/docker"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
)

var tmpDockerfile = "Dockerfile"

var DockerfileYmlOrders = []string{
	"build/dockerfile.default.yml",
	"build/dockerfile.yml",
}

func CommandsForShipping(p *Project, push bool) (commands []*exec.Cmd) {
	dockerfile := &docker.Dockerfile{}

	hasDockerfileYaml := false

	for _, dockerfileYml := range DockerfileYmlOrders {
		if isPathExist(dockerfileYml) {
			hasDockerfileYaml = true
			mayReadFileAndUnmarshal(dockerfileYml, dockerfile)
		}
	}

	if dockerfile.Image == "" {
		dockerfile.Image = "${PROFZONE_DOCKER_REGISTRY}/${PROJECT_OWNER}/${PROJECT_NAME}:${PROJECT_VERSION}"
	}

	if hasDockerfileYaml {
		dockerfile.AddEnv(EnvKeyProjectVersion, p.Version.String())
		dockerfile.AddEnv(EnvKeyProjectOwner, p.Owner)
		dockerfile.AddEnv(EnvKeyProjectGroup, p.Group)
		dockerfile.AddEnv(EnvKeyProjectName, p.Name)
		dockerfile.AddEnv(EnvKeyProjectFeature, p.Feature)

		ioutil.WriteFile(tmpDockerfile, []byte(dockerfile.String()), os.ModePerm)
	}

	commands = append(commands, p.Command("docker", "build", "-f", tmpDockerfile, "-t", dockerfile.Image, "."))
	if push {
		commands = append(commands, p.Command("docker", "push", dockerfile.Image))
	}
	return
}

func isPathExist(path string) bool {
	f, _ := os.Stat(path)
	return f != nil
}

func mayReadFileAndUnmarshal(file string, v interface{}) {
	bytes, errForRead := ioutil.ReadFile(file)
	if errForRead != nil {
		panic(errForRead)
	}
	err := yaml.Unmarshal(bytes, v)
	if err != nil {
		panic(err)
	}
}
