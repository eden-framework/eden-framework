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
	"github.com/profzone/eden-framework/internal/project"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ciDeployCmdEnv               string
	ciDeployCmdConfigFile        string
	ciDeployCmdDeployConfigFile  string
	ciDeployCmdServiceConfigFile string
)

// ciShipCmd represents the ciShip command
var ciDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "ci ship a project as a image",
	Run: func(cmd *cobra.Command, args []string) {
		err := project.ProcessDeployment(currentProject, ciDeployCmdEnv, ciDeployCmdConfigFile, ciDeployCmdDeployConfigFile, ciDeployCmdServiceConfigFile)
		if err != nil {
			logrus.Panic(err)
		}
	},
}

func init() {
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdEnv, "env", "e", "", "deploy environment name")
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdConfigFile, "config", "c", "", "kubeconfig file path")
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdDeployConfigFile, "deploy", "d", "./build/deploy.default.yml", "deploy yaml file path")
	ciDeployCmd.Flags().StringVarP(&ciDeployCmdServiceConfigFile, "service", "s", "./build/service.default.yml", "service yaml file path")
	ciCmd.AddCommand(ciDeployCmd)
}
