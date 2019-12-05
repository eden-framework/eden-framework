package generator

import (
	"github.com/profzone/eden-framework/internal/generator/files"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"path"
)

const (
	Image     = "${PROFZONE_DOCKER_REGISTRY}/${PROJECT_GROUP}/${PROJECT_NAME}:${PROJECT_VERSION}"
	FromImage = "${PROFZONE_DOCKER_REGISTRY}/profzone/golang:runtime"
)

type DockerGenerator struct {
	ServiceName string
}

func NewDockerGenerator(serviceName string) *DockerGenerator {
	return &DockerGenerator{
		ServiceName: serviceName,
	}
}

func (d *DockerGenerator) Load(path string) {
}

func (d *DockerGenerator) Pick() {
}

func (d *DockerGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}

	dockerFile := &files.Dockerfile{
		From:  FromImage,
		Image: Image,
	}
	dockerFile = dockerFile.AddEnv("GOENV", "DEV")

	dockerFile = dockerFile.WithWorkDir("/go/bin")
	dockerFile = dockerFile.WithCmd("./"+d.ServiceName, "-d=false", "-m=false")
	dockerFile = dockerFile.WithExpose("80")

	dockerFile = dockerFile.AddContent("./configs", "./configs")
	dockerFile = dockerFile.AddContent("./build/"+d.ServiceName, "./")
	dockerFile = dockerFile.AddContent("./profzone.yml", "./")
	dockerFile = dockerFile.AddContent("./api/api.json", "./")

	content, err := yaml.Marshal(dockerFile)
	if err != nil {
		logrus.Panic(err)
	}

	outputs.Add(path.Join(outputPath, "dockerfile.default.yml"), string(content))
	return outputs
}
