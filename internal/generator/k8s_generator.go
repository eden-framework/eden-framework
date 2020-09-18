package generator

import (
	"bytes"
	"fmt"
	"github.com/profzone/eden-framework/internal"
	"github.com/profzone/eden-framework/internal/generator/files"
	"github.com/profzone/eden-framework/internal/project"
	"github.com/profzone/envconfig"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	"path"
)

const (
	EnvVarKeyDeploymentUID = "DEPLOYMENT_UID"
)

type K8sGenerator struct {
	ServiceName string
	EnvVars     []envconfig.EnvVar
}

func NewK8sGenerator(serviceName string, envVars []envconfig.EnvVar) *K8sGenerator {
	return &K8sGenerator{
		ServiceName: serviceName,
		EnvVars:     envVars,
	}
}

func (d *K8sGenerator) Load(path string) {
}

func (d *K8sGenerator) Pick() {
}

func (d *K8sGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}

	// Deployment config file
	deploymentConfig := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      files.EnvVar(internal.EnvVarKeyProjectName),
			Namespace: files.EnvVar(internal.EnvVarKeyProjectGroup),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int2Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"workload.user.cattle.io/workloadselector": files.EnvVar(internal.EnvVarKeyProjectSelector),
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"workload.user.cattle.io/workloadselector": files.EnvVar(internal.EnvVarKeyProjectSelector),
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: files.EnvVar(internal.EnvVarKeyProjectName),
							Image: fmt.Sprintf("%s/%s/%s:%s",
								files.EnvVar(project.DOCKER_REGISTRY_KEY),
								files.EnvVar(internal.EnvVarKeyProjectGroup),
								files.EnvVar(internal.EnvVarKeyProjectName),
								files.EnvVar(internal.EnvVarKeyProjectRef),
							),
							ImagePullPolicy: apiv1.PullAlways,
							Ports:           []apiv1.ContainerPort{},
						},
					},
					RestartPolicy: apiv1.RestartPolicyAlways,
				},
			},
		},
	}

	// Service config file
	serviceConfig := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      files.EnvVar(internal.EnvVarKeyProjectName),
			Namespace: files.EnvVar(internal.EnvVarKeyProjectGroup),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1beta2",
					Kind:       "Deployment",
					Name:       files.EnvVar(internal.EnvVarKeyProjectName),
					UID:        types.UID(files.EnvVar(EnvVarKeyDeploymentUID)),
				},
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"workload.user.cattle.io/workloadselector": files.EnvVar(internal.EnvVarKeyProjectSelector),
			},
			Ports: []apiv1.ServicePort{},
		},
	}

	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{
		Yaml:   true,
		Pretty: true,
		Strict: true,
	})
	deploymentBuffer := bytes.NewBuffer([]byte{})
	err := serializer.Encode(deploymentConfig, deploymentBuffer)
	if err != nil {
		logrus.Panic(err)
	}
	serviceBuffer := bytes.NewBuffer([]byte{})
	err = serializer.Encode(serviceConfig, serviceBuffer)
	if err != nil {
		logrus.Panic(err)
	}

	outputs.Add(path.Join(outputPath, "build/deploy.default.yml"), deploymentBuffer.String())
	outputs.Add(path.Join(outputPath, "build/service.default.yml"), serviceBuffer.String())
	return outputs
}
