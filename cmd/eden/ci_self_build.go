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
	"fmt"
	"gitee.com/eden-framework/eden-framework/internal/project"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ciSelfBuildCmd represents the selfBuild command
var ciSelfBuildCmd = &cobra.Command{
	Use:   "self-build",
	Short: "build eden tool-chain",
	Run: func(cmd *cobra.Command, args []string) {
		currentProject.Run(currentProject.Command(
			"docker",
			"build",
			"-t registry.cn-hangzhou.aliyuncs.com/profzone/golang-onbuild:1.14",
			"-f scripts/onbuild.Dockerfile",
			"."))
		currentProject.Run(currentProject.Command(
			"docker",
			"build",
			"-t registry.cn-hangzhou.aliyuncs.com/profzone/golang-onship:1.14",
			"-f scripts/onship.Dockerfile",
			"."))
		currentProject.Run(currentProject.Command(
			"docker",
			"build",
			"-t registry.cn-hangzhou.aliyuncs.com/profzone/golang-ondeploy:2.4.3",
			"-f scripts/ondeploy.Dockerfile",
			"."))

		processor := viper.GetString("SHIPPING_PROCESSOR")
		typ, err := project.ParseShippingProcessorTypeFromString(processor)
		if err != nil {
			panic(fmt.Sprintf("cannot parse shipping processor type from env: SHIPPING_PROCESSOR=%s", processor))
		}
		shipping := project.NewShippingProcessor(typ)
		currentProject.Run(shipping.Login(currentProject)...)
		currentProject.Run(shipping.Push(currentProject, "registry.cn-hangzhou.aliyuncs.com/profzone/golang-onbuild:1.14")...)
		currentProject.Run(shipping.Push(currentProject, "registry.cn-hangzhou.aliyuncs.com/profzone/golang-onship:1.14")...)
		currentProject.Run(shipping.Push(currentProject, "registry.cn-hangzhou.aliyuncs.com/profzone/golang-ondeploy:2.4.3")...)
	},
}

func init() {
	ciCmd.AddCommand(ciSelfBuildCmd)
}
