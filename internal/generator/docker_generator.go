package generator

import (
	"github.com/profzone/eden-framework/internal/generator/files"
	"github.com/profzone/envconfig"
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
	EnvVars     []envconfig.EnvVar
}

func NewDockerGenerator(serviceName string, envVars []envconfig.EnvVar) *DockerGenerator {
	return &DockerGenerator{
		ServiceName: serviceName,
		EnvVars:     envVars,
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

	for _, envVar := range d.EnvVars {
		strValue := envVar.GetValue(false)
		dockerFile = dockerFile.AddEnv(envVar.Key, strValue)
	}

	dockerFile = dockerFile.WithWorkDir("/go/bin")
	dockerFile = dockerFile.WithCmd("./"+d.ServiceName, "-d=false", "-m=false")
	dockerFile = dockerFile.WithExpose("80")

	dockerFile = dockerFile.AddContent("./build/configs", "./configs")
	dockerFile = dockerFile.AddContent("./build/"+d.ServiceName, "./")
	dockerFile = dockerFile.AddContent("./profzone.yml", "./")
	dockerFile = dockerFile.AddContent("./openapi.json", "./")

	content, err := yaml.Marshal(dockerFile)
	if err != nil {
		logrus.Panic(err)
	}

	configDefaultFile := files.NewConfigDefaultFile(d.EnvVars)

	outputs.Add(path.Join(outputPath, "build/dockerfile.default.yml"), string(content))
	outputs.Add(path.Join(outputPath, "build/configs/default.yml"), configDefaultFile.String())
	return outputs
}
