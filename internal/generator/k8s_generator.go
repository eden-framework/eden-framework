package generator

import (
	"bytes"
	"fmt"
	"github.com/eden-framework/courier/transport_grpc"
	"github.com/eden-framework/courier/transport_http"
	"github.com/eden-framework/eden-framework/internal/generator/files"
	"github.com/eden-framework/eden-framework/internal/project"
	"github.com/eden-framework/pointer"
	"github.com/eden-framework/reflectx"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"path"
	"reflect"
	"strings"
)

type K8sGenerator struct {
	config []interface{}
}

func NewK8sGenerator(config []interface{}) *K8sGenerator {
	return &K8sGenerator{
		config: config,
	}
}

func (d *K8sGenerator) Load(cwd string) {
}

func (d *K8sGenerator) Pick() {
}

func (d *K8sGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}

	serverPorts := findServerPorts(d.config)

	// Deployment config file
	deploymentConfig := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      files.EnvVar(project.EnvKeyProjectName),
			Namespace: files.EnvVar(project.EnvKeyProjectGroup),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"workload.user.cattle.io/workloadselector": files.EnvVar(project.EnvKeyProjectSelector),
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"workload.user.cattle.io/workloadselector": files.EnvVar(project.EnvKeyProjectSelector),
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: files.EnvVar(project.EnvKeyProjectName),
							Image: fmt.Sprintf("%s/%s/%s:%s",
								files.EnvVar(project.EnvKeyDockerRegistryKey),
								files.EnvVar(project.EnvKeyProjectGroup),
								files.EnvVar(project.EnvKeyProjectName),
								files.EnvVar(project.EnvKeyProjectVersion),
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
			Name:      files.EnvVar(project.EnvKeyProjectName),
			Namespace: files.EnvVar(project.EnvKeyProjectGroup),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1beta2",
					Kind:       "Deployment",
					Name:       files.EnvVar(project.EnvKeyProjectName),
					UID:        types.UID(files.EnvVar(project.EnvKeyDeploymentUID)),
				},
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"workload.user.cattle.io/workloadselector": files.EnvVar(project.EnvKeyProjectSelector),
			},
			Ports: []apiv1.ServicePort{},
		},
	}

	for _, port := range serverPorts {
		deploymentConfig.Spec.Template.Spec.Containers[0].Ports = append(deploymentConfig.Spec.Template.Spec.Containers[0].Ports, apiv1.ContainerPort{
			Name:          fmt.Sprintf("tcp%d", port),
			ContainerPort: port,
			Protocol:      apiv1.ProtocolTCP,
		})
		serviceConfig.Spec.Ports = append(serviceConfig.Spec.Ports, apiv1.ServicePort{
			Name:       fmt.Sprintf("tcp%d", port),
			Protocol:   apiv1.ProtocolTCP,
			Port:       port,
			TargetPort: intstr.FromInt(int(port)),
		})
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

	outputs.Add(path.Join(outputPath, "./deploy.default.yml"), deploymentBuffer.String())
	outputs.Add(path.Join(outputPath, "./service.default.yml"), serviceBuffer.String())
	return outputs
}

func (d *K8sGenerator) Finally() {

}

func findServerPorts(config []interface{}) (ports []int32) {
	for _, c := range config {
		val := reflect.ValueOf(c)

		exampleNames := exampleServerFullTypeName()
		EachFieldValue(val.Elem(), func(field reflect.Value) bool {
			fullName := reflectx.FullTypeName(reflectx.FromRType(field.Type()))
			for _, name := range exampleNames {
				if strings.HasSuffix(fullName, name) {
					field = reflectx.Indirect(field)
					portVal := field.FieldByName("Port")
					if !reflectx.IsEmptyValue(portVal) {
						ports = append(ports, int32(portVal.Int()))
					}
				}
			}

			return true
		})
	}
	return
}

func exampleServerFullTypeName() (names []string) {
	exampleHttpServer := transport_http.ServeHTTP{}
	exampleGrpcServer := transport_grpc.ServeGRPC{}

	httpTyp := reflect.TypeOf(exampleHttpServer)
	names = append(names, httpTyp.PkgPath()+"."+httpTyp.Name())
	grpcTyp := reflect.TypeOf(exampleGrpcServer)
	names = append(names, grpcTyp.PkgPath()+"."+grpcTyp.Name())
	return
}

func EachFieldValue(val reflect.Value, walker func(value reflect.Value) bool) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !walker(field) {
			break
		}
	}
}
