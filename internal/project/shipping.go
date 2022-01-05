package project

import (
	"fmt"
	"gitee.com/eden-framework/eden-framework/internal/docker"
	"github.com/spf13/viper"
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
	if !hasDockerfileYaml {
		panic("there has no dockerfile.yml file in project workspace")
	}

	if dockerfile.Image == "" {
		dockerfile.Image = "${PROFZONE_DOCKER_REGISTRY}/${PROJECT_GROUP}/${PROJECT_NAME}:${PROJECT_VERSION}"
	}

	dockerfile.AddEnv(EnvKeyProjectVersion, p.Version.String())
	dockerfile.AddEnv(EnvKeyProjectOwner, p.Owner)
	dockerfile.AddEnv(EnvKeyProjectGroup, p.Group)
	dockerfile.AddEnv(EnvKeyProjectName, p.Name)
	dockerfile.AddEnv(EnvKeyProjectFeature, p.Feature)

	dockerfileContent := dockerfile.String()
	fmt.Println(dockerfileContent)
	ioutil.WriteFile(tmpDockerfile, []byte(dockerfileContent), os.ModePerm)

	commands = append(commands, p.Command("docker", "build", "-f", tmpDockerfile, "-t", dockerfile.Image, "."))
	if push {
		processor := viper.GetString("SHIPPING_PROCESSOR")
		typ, err := ParseShippingProcessorTypeFromString(processor)
		if err != nil {
			panic(fmt.Sprintf("cannot parse shipping processor type from env: SHIPPING_PROCESSOR=%s", processor))
		}
		shipping := NewShippingProcessor(typ)
		commands = append(commands, shipping.Login(p)...)
		commands = append(commands, shipping.Push(p, dockerfile.Image)...)
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
