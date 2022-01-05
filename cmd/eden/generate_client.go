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
	"fmt"
	"gitee.com/eden-framework/eden-framework/internal/generator"
	"github.com/spf13/cobra"
	"os"
)

var (
	clientCmdApiPath, clientCmdOutputPath, clientCmdPackageName string
)

// clientCmd represents the files command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long:  fmt.Sprintf("%s\ngenerate files", CommandHelpHeader),
	Run: func(cmd *cobra.Command, args []string) {
		if clientCmdOutputPath == "" {
			clientCmdOutputPath, _ = os.Getwd()
		}
		gen := generator.NewClientGenerator(clientCmdPackageName)
		generator.Generate(gen, clientCmdApiPath, clientCmdOutputPath)
	},
}

func init() {
	generateCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&clientCmdApiPath, "api-path", "p", "", "eden generate api --api-path=/go/src/eden-server/api.json")
	clientCmd.Flags().StringVarP(&clientCmdOutputPath, "output-path", "o", "", "eden generate api --output-path=/go/src/eden-server/")
	clientCmd.Flags().StringVarP(&clientCmdPackageName, "service-name", "n", "", "eden generate api --service-name=account")
}
