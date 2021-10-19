package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eden-framework/eden-framework/internal/k8s"
	str "github.com/eden-framework/strings"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

func ProcessDeployment(p *Project, env, deployConfig, serviceConfig string) error {
	if env == "" {
		return errors.New("deployment must specify a environment name")
	}
	var envVars map[string]string
	fmt.Printf("CURRENT env: %s\n", env)
	if strings.ToLower(env) != "prod" {
		fmt.Printf("strings.ToLower(env): %s\n", strings.ToLower(env))
		envVars = LoadEnv(env, p.Feature)
	}

	kubeConfigKey := str.ToUpperSnakeCase("KubeConfig" + env)
	kubeConfig := viper.GetString(kubeConfigKey)
	if len(kubeConfig) == 0 {
		panic("cannot find kube config file path from .eden.yaml, the key is " + kubeConfigKey)
	}

	ctx, _ := context.WithCancel(context.Background())
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Get Deployment
	deployment, patch, err := k8s.MakeDeployment(deployConfig)
	if err != nil {
		return err
	}

	for key, val := range envVars {
		deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, apiv1.EnvVar{
			Name:  key,
			Value: val,
		})
	}
	patch = patchEnvVars(patch, envVars)

	deploymentsClient := clientSet.AppsV1().Deployments(deployment.Namespace)
	_, err = deploymentsClient.Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		// Create Deployment
		fmt.Println("Creating deployment...")
		deployment, err = deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Created deployment %s.\n", deployment.GetObjectMeta().GetName())
	} else {
		// Patch Deployment
		fmt.Println("Updating deployment...")
		data, err := json.Marshal(patch)
		if err != nil {
			return err
		}
		deployment, err = deploymentsClient.Patch(ctx, deployment.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Updated deployment %s.\n", deployment.GetObjectMeta().GetName())
	}

	SetEnv(EnvKeyDeploymentUID, string(deployment.UID))

	// Get Service
	service, patch, err := k8s.MakeService(serviceConfig)
	if err != nil {
		return err
	}
	servicesClient := clientSet.CoreV1().Services(service.Namespace)
	_, err = servicesClient.Get(ctx, service.Name, metav1.GetOptions{})
	if err != nil {
		// Create Service
		fmt.Println("Creating service...")
		service, err = servicesClient.Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Created service %s.\n", service.GetObjectMeta().GetName())
	} else {
		// Patch Service
		fmt.Println("Updating service...")
		data, err := json.Marshal(patch)
		if err != nil {
			return err
		}
		service, err = servicesClient.Patch(ctx, service.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Updated service %s.\n", service.GetObjectMeta().GetName())
	}
	return nil
}

func LoadEnv(envName string, feature string) map[string]string {
	defaultEnv := loadEnvFromFiles("default", feature)
	if envName != "" {
		extendEnv := loadEnvFromFiles(envName, feature)
		defaultEnv = mergeEnvVars(defaultEnv, extendEnv)
	}

	return defaultEnv
}

func loadEnvFromFiles(envName string, feature string) map[string]string {
	defaultEnv := loadEnvFromFile(envName)
	if feature != "" {
		extendEnv := loadEnvFromFile(envName + "-" + feature)
		defaultEnv = mergeEnvVars(defaultEnv, extendEnv)
	}

	return defaultEnv
}

func loadEnvFromFile(envName string) map[string]string {
	filename := "build/configs/" + strings.ToLower(envName) + ".yml"
	logrus.Infof("try to load env vars from %s ...", color.GreenString(filename))
	envFileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	var envVars map[string]string
	err = yaml.Unmarshal(envFileContent, &envVars)
	if err != nil {
		panic(err)
	}
	for key, value := range envVars {
		SetEnv(key, value)
	}
	return envVars
}

func mergeEnvVars(self map[string]string, source map[string]string) map[string]string {
	for key, val := range source {
		self[key] = val
	}
	return self
}

func patchEnvVars(patch map[string]interface{}, envVars map[string]string) map[string]interface{} {
	envs := make([]map[string]string, 0)
	for key, val := range envVars {
		envs = append(envs, map[string]string{
			"name":  key,
			"value": val,
		})
	}
	patch["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["env"] = envs

	return patch
}
