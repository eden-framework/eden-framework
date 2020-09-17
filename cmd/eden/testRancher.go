/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
)

// testRancherCmd represents the rancher command
var testRancherCmd = &cobra.Command{
	Use:   "testRancher",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		//currentProject.Run(project.CommandsForShipping(currentProject, shipFlagPush)...)
		//if ciDeployCmdProjectName == "" {
		//	panic("the flag --project must not be empty")
		//}
		//clientbase.Debug = false
		//baseClient, err := client.NewClient(&clientbase.ClientOpts{
		//	URL:       "https://centos:8443/v3",
		//	AccessKey: "token-tj5xv",
		//	SecretKey: "w6jc65w5j9bmczt5cds2t7f2wnn7vpgltwjdpblkbhjjmxvk5r6rvw",
		//	Timeout:   10 * time.Second,
		//	Insecure:  true,
		//})
		//if err != nil {
		//	logrus.Panicf("NewClient err: %v", err)
		//}
		//
		//var project client.Project
		//projects, err := baseClient.Project.List(&types.ListOpts{
		//	Filters: map[string]interface{}{
		//		"name": ciDeployCmdProjectName,
		//	},
		//})
		//if err != nil {
		//	logrus.Panicf("Project.List err: %v", err)
		//}
		//
		//if len(projects.Data) == 0 {
		//	p, err := baseClient.Project.Create(&client.Project{
		//		Name: ciDeployCmdProjectName,
		//	})
		//	if err != nil {
		//		logrus.Panicf("Project.Create err: %v", err)
		//	}
		//	project = *p
		//} else {
		//	project = projects.Data[0]
		//}
		//
		//projectClient, err := pClient.NewClient(&clientbase.ClientOpts{
		//	URL:       "https://centos:8443/v3/projects/" + project.ID,
		//	AccessKey: "token-tj5xv",
		//	SecretKey: "w6jc65w5j9bmczt5cds2t7f2wnn7vpgltwjdpblkbhjjmxvk5r6rvw",
		//	Timeout:   10 * time.Second,
		//	Insecure:  true,
		//})
		//if err != nil {
		//	logrus.Panicf("NewClient err: %v", err)
		//}
		//
		//deployments, err := projectClient.Deployment.List(&types.ListOpts{
		//	Filters: map[string]interface{}{
		//		"name": ciDeployCmdDeployName,
		//	},
		//})
		//if err != nil {
		//	logrus.Panicf("Deployment.List err: %v", err)
		//}
		//if len(deployments.Data) == 0 {
		//	// create
		//} else {
		//	// update
		//	deployment := deployments.Data[0]
		//	_, err = projectClient.Deployment.Update(&deployment, map[string]interface{}{})
		//}
	},
}

//func GetProjectByName(name string, collection *client.ProjectCollection) *client.Project {
//	for _, p := range collection.Data {
//		if p.Name == name {
//			return &p
//		}
//	}
//	return nil
//}
//
//func GetDeploymentByName(name string, collection *pClient.DeploymentCollection) *pClient.Deployment {
//	for _, d := range collection.Data {
//		if d.Name == name {
//			return &d
//		}
//	}
//	return nil
//}
//
//func init() {
//	rootCmd.AddCommand(testRancherCmd)
//}
