/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/profzone/eden-framework/internal/k8s"
	"github.com/profzone/eden-framework/internal/project"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ciDeployCmdNamespace         string
	ciDeployCmdConfigFile        string
	ciDeployCmdDeployConfigFile  string
	ciDeployCmdServiceConfigFile string
)

const (
	DeploymentUIDEnvVarKey = "DEPLOYMENT_UID"
)

// ciShipCmd represents the ciShip command
var ciDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "ci ship a project as a image",
	Run: func(cmd *cobra.Command, args []string) {
		currentProject.SetEnviron()

		ctx, _ := context.WithCancel(context.Background())
		config, err := clientcmd.BuildConfigFromFlags("", ciDeployCmdConfigFile)
		if err != nil {
			panic(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		// Get Deployment
		deploymentsClient := clientset.AppsV1().Deployments(ciDeployCmdNamespace)
		deployment, patch, err := k8s.MakeDeployment(ciDeployCmdDeployConfigFile)
		if err != nil {
			panic(err)
		}
		_, err = deploymentsClient.Get(ctx, deployment.Name, metav1.GetOptions{})
		if err != nil {
			// Create Deployment
			fmt.Println("Creating deployment...")
			deployment, err = deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Created deployment %s.\n", deployment.GetObjectMeta().GetName())
		} else {
			// Patch Deployment
			fmt.Println("Updating deployment...")
			data, err := json.Marshal(patch)
			if err != nil {
				panic(err)
			}
			deployment, err = deploymentsClient.Patch(ctx, deployment.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Updated deployment %s.\n", deployment.GetObjectMeta().GetName())
		}

		project.SetEnv(DeploymentUIDEnvVarKey, string(deployment.UID))

		// Get Service
		servicesClient := clientset.CoreV1().Services(ciDeployCmdNamespace)
		service, patch, err := k8s.MakeService(ciDeployCmdServiceConfigFile)
		if err != nil {
			panic(err)
		}
		_, err = servicesClient.Get(ctx, service.Name, metav1.GetOptions{})
		if err != nil {
			// Create Service
			fmt.Println("Creating service...")
			service, err = servicesClient.Create(ctx, service, metav1.CreateOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Created service %s.\n", service.GetObjectMeta().GetName())
		} else {
			// Patch Service
			fmt.Println("Updating service...")
			data, err := json.Marshal(patch)
			if err != nil {
				panic(err)
			}
			service, err = servicesClient.Patch(ctx, service.Name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Updated service %s.\n", service.GetObjectMeta().GetName())
		}
	},
}

func init() {
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdConfigFile, "config", "c", "", "kubeconfig file path")
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdDeployConfigFile, "deploy", "d", "./build/deploy.yml", "deploy yaml file path")
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdServiceConfigFile, "service", "s", "./build/service.yml", "service yaml file path")
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdNamespace, "namespace", "n", "default", "eden ci deploy --namespace=default")
	ciCmd.AddCommand(ciDeployCmd)
}
