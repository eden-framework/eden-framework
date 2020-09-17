package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/profzone/eden-framework/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DeploymentUIDEnvVarKey = "DEPLOYMENT_UID"
)

func ProcessDeployment(kubeConfig, deployConfig, serviceConfig string) error {
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

	SetEnv(DeploymentUIDEnvVarKey, string(deployment.UID))

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
