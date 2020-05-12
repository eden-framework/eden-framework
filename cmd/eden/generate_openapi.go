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
	"github.com/profzone/eden-framework/internal/generator"
	"github.com/profzone/eden-framework/internal/generator/scanner"
	"os"

	"github.com/spf13/cobra"
)

var openApiCmdInputPath, openApiCmdOutputPath string

// apiCmd represents the api command
var openApiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long:  fmt.Sprintf("%s\ngenerate api doc", CommandHelpHeader),
	Run: func(cmd *cobra.Command, args []string) {
		if openApiCmdInputPath == "" {
			openApiCmdInputPath, _ = os.Getwd()
		}
		if openApiCmdOutputPath == "" {
			openApiCmdOutputPath, _ = os.Getwd()
		}
		enumScanner := scanner.NewEnumScanner()
		modelScanner := scanner.NewModelScanner(enumScanner)
		operatorScanner := scanner.NewOperatorScanner(modelScanner)
		gen := generator.NewOpenApiGenerator(operatorScanner, modelScanner, enumScanner)

		modelScanner.Api = &gen.Api
		operatorScanner.Api = &gen.Api

		generator.Generate(gen, openApiCmdInputPath, openApiCmdOutputPath)
	},
}

func init() {
	generateCmd.AddCommand(openApiCmd)
	openApiCmd.Flags().StringVarP(&openApiCmdInputPath, "input-path", "i", "", "eden generate api --input-path=/go/src/eden-server")
	openApiCmd.Flags().StringVarP(&openApiCmdOutputPath, "output-path", "o", "", "eden generate api --output-path=/go/src/eden-server")
}
