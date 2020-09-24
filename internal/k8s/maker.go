package k8s

import (
	"bytes"
	"fmt"
	"github.com/eden-framework/eden-framework/pkg/executil"
	"github.com/pkg/errors"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
)

func MakeDeployment(pathToYaml string) (*appsv1.Deployment, map[string]interface{}, error) {
	manifest, err := PathToOSFile(pathToYaml)
	if err != nil {
		return nil, nil, err
	}
	manifestData, err := ioutil.ReadAll(manifest)
	if err != nil {
		return nil, nil, err
	}

	envVars := executil.EnvVars{}
	envVars.LoadFromEnviron()
	manifestData = []byte(envVars.Parse(string(manifestData)))
	reader := bytes.NewReader(manifestData)

	deployment := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(reader, 100).Decode(&deployment); err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("failed to decode file %s into Deployment", pathToYaml))
	}

	reader.Reset(manifestData)

	patch := make(map[string]interface{})
	if err := yaml.NewYAMLOrJSONDecoder(reader, 100).Decode(&patch); err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("failed to decode file %s into patch", pathToYaml))
	}

	return &deployment, patch, nil
}

func MakeService(pathToYaml string) (*apiv1.Service, map[string]interface{}, error) {
	manifest, err := PathToOSFile(pathToYaml)
	if err != nil {
		return nil, nil, err
	}
	manifestData, err := ioutil.ReadAll(manifest)
	if err != nil {
		return nil, nil, err
	}

	envVars := executil.EnvVars{}
	envVars.LoadFromEnviron()
	manifestData = []byte(envVars.Parse(string(manifestData)))
	reader := bytes.NewReader(manifestData)

	resource := apiv1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(reader, 100).Decode(&resource); err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("failed to decode file %s into Service", pathToYaml))
	}
	reader.Reset(manifestData)

	patch := make(map[string]interface{})
	if err := yaml.NewYAMLOrJSONDecoder(reader, 100).Decode(&patch); err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("failed to decode file %s into patch", pathToYaml))
	}

	return &resource, patch, nil
}

func PathToOSFile(relativPath string) (*os.File, error) {
	path, err := filepath.Abs(relativPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed generate absolut file path of %s", relativPath))
	}

	manifest, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to open file %s", path))
	}

	return manifest, nil
}
