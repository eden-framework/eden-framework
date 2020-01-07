package project

import (
	"github.com/fatih/color"
	"github.com/profzone/eden-framework/internal/docker"
	"github.com/profzone/eden-framework/pkg/ptr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	EnvVarRef              = "PROJECT_REF"
	EnvVarBuildRef         = "DRONE_COMMIT_SHA"
	EnvVarBuildBranch      = "DRONE_COMMIT_REF"
	EnvVarRancherEnv       = "RANCHER_ENVIRONMENT"
	EnvVarRancherUrl       = "RANCHER_URL"
	EnvVarRancherAccessKey = "RANCHER_ACCESS_KEY"
	EnvVarRancherSecretKey = "RANCHER_SECRET_KEY"
)

var (
	EnvValRancherUrl       = "http://rancher.profzone.net:38080"
	EnvValRancherAccessKey = "744E0D8EF311C269FED1"
	EnvValRancherSecretKey = "yBXzp7jdaaRqCtL92TJSRbekxzYr8x7Xr2r5rq11"
)

var tmpDockerfile = "Dockerfile"

var DockerfileYmlOrders = []string{
	"dockerfile.default.yml",
	"dockerfile.yml",
}

var (
	CIWorkingDirectory  = "/drone/workspace"
	CIGolangRootPath    = "/go/src/"
	COGolangPackageName = "github.com/"
)

func CommandForDeploy(p *Project, deployEnv string) (command *exec.Cmd) {
	SetEnv(EnvVarRancherEnv, deployEnv)
	if viper.GetString("RANCHER_URL") == "" {
		SetEnv(EnvVarRancherUrl, EnvValRancherUrl)
	} else {
		SetEnv(EnvVarRancherUrl, viper.GetString("RANCHER_URL"))
	}
	if viper.GetString("RANCHER_ACCESS_KEY") == "" {
		SetEnv(EnvVarRancherAccessKey, EnvValRancherAccessKey)
	} else {
		SetEnv(EnvVarRancherAccessKey, viper.GetString("RANCHER_ACCESS_KEY"))
	}
	if viper.GetString("RANCHER_SECRET_KEY") == "" {
		SetEnv(EnvVarRancherSecretKey, EnvValRancherSecretKey)
	} else {
		SetEnv(EnvVarRancherSecretKey, viper.GetString("RANCHER_SECRET_KEY"))
	}
	stackName := p.Group

	if p.Feature != "" {
		stackName = stackName + "--" + p.Feature
	}

	LoadEnv(deployEnv, p.Feature)

	writeMemoryLimit(p.Name)

	rancherUp := []string{
		"rancher",
		"up",
		"-d",
	}

	_, err := os.Stat("/usr/local/bin/rancher-env.sh")
	if err == nil {
		rancherUp = append([]string{"rancher-env.sh"}, rancherUp...)
	}

	dockerComposeFiles := []string{
		"docker-compose.initial.yml",
		"docker-compose.default.yml",
		"docker-compose.yml",
	}

	for _, dockerComposeFile := range dockerComposeFiles {
		if isPathExist(dockerComposeFile) {
			rancherUp = append(rancherUp, "-f", dockerComposeFile)
		}
	}

	if p.Feature != "" {
		p.Version.Prefix = p.Feature
	}

	rancherUp = append(rancherUp, "--stack", stackName, "--pull", "--force-upgrade", "--confirm-upgrade")

	return p.Command(rancherUp...)
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
		dockerfile.Image = "${PROFZONE_DOCKER_REGISTRY}/${PROJECT_GROUP}/${PROJECT_NAME}:${PROJECT_VERSION}"
	}

	if hasDockerfileYaml {
		p.SetEnviron()
		dockerfile = dockerfile.AddEnv(EnvVarRef, p.Version.String()+"-"+os.Getenv(EnvVarBuildRef))

		dockerfile.AddEnv("PROJECT_GROUP", p.Group)
		dockerfile.AddEnv("PROJECT_NAME", p.Name)
		dockerfile.AddEnv("PROJECT_FEATURE", p.Feature)

		ioutil.WriteFile(tmpDockerfile, []byte(dockerfile.String()), os.ModePerm)
	}

	if p.Feature != "" {
		p.Version.Prefix = p.Feature
	}

	commands = append(commands, p.Command("docker", "build", "-f", tmpDockerfile, "-t", dockerfile.Image, "."))
	if push {
		commands = append(commands, p.Command("docker", "push", dockerfile.Image))
	}
	return
}

func writeMemoryLimit(serviceName string) {
	compose := docker.NewDockerCompose()

	s := docker.NewService("busybox:latest")
	s.MemLimit = ptr.Int64(1073741824)

	compose = compose.AddService(serviceName, s)
	data, _ := yaml.Marshal(compose)

	ioutil.WriteFile("docker-compose.initial.yml", data, os.ModePerm)
}

func isPathExist(path string) bool {
	f, _ := os.Stat(path)
	return f != nil
}

func LoadEnv(envName string, feature string) {
	loadEnvFromFiles("default", feature)
	if envName != "" {
		loadEnvFromFiles(envName, feature)
	}
}

func loadEnvFromFiles(envName string, feature string) {
	loadEnvFromFile(envName)
	if feature != "" {
		loadEnvFromFile(envName + "-" + feature)
	}
}

func loadEnvFromFile(envName string) {
	filename := "config/" + strings.ToLower(envName) + ".yml"
	logrus.Infof("try to load env vars from %s ...\n", color.GreenString(filename))
	envFileContent, err := ioutil.ReadFile(filename)
	if err == nil {
		var envVars map[string]string
		err := yaml.Unmarshal([]byte(envFileContent), &envVars)
		if err != nil {
			panic(err)
		}
		for key, value := range envVars {
			SetEnv(key, value)
		}
	}
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
