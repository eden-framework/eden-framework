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
	"os"
	"profzone/eden-framework/internal/generator"
	"profzone/eden-framework/internal/generator/scanner"

	"github.com/spf13/cobra"
)

var apiCmdCWD string

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long:  fmt.Sprintf("%s\ngenerate api doc", CommandHelpHeader),
	Run: func(cmd *cobra.Command, args []string) {
		if apiCmdCWD == "" {
			apiCmdCWD, _ = os.Getwd()
		}
		modelScanner := scanner.NewModelScanner()
		operatorScanner := scanner.NewOperatorScanner(modelScanner)
		gen := generator.NewApiGenerator(operatorScanner, modelScanner)

		modelScanner.Api = &gen.Api
		operatorScanner.Api = &gen.Api

		generator.Generate(gen, apiCmdCWD)
	},
}

func init() {
	generateCmd.AddCommand(apiCmd)
	apiCmd.Flags().StringVarP(&apiCmdCWD, "input-path", "i", "", "eden generate api --input-path=/go/src/eden-server")
}