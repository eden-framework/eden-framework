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
	"github.com/spf13/cobra"
	"os"
)

var enumCmdInputPath, enumCmdOutputPath, enumCmdTypeName string

// enumCmd represents the enum command
var enumCmd = &cobra.Command{
	Use:   "enum",
	Short: "A brief description of your command",
	Long:  fmt.Sprintf("%s\ngenerate enum", CommandHelpHeader),
	Run: func(cmd *cobra.Command, args []string) {
		if enumCmdInputPath == "" {
			enumCmdInputPath, _ = os.Getwd()
		}
		if enumCmdOutputPath == "" {
			enumCmdOutputPath, _ = os.Getwd()
		}
		enumScanner := scanner.NewEnumScanner()
		gen := generator.NewEnumGenerator(enumScanner, enumCmdTypeName)

		generator.Generate(gen, enumCmdInputPath, enumCmdOutputPath)
	},
}

func init() {
	generateCmd.AddCommand(enumCmd)
	enumCmd.Flags().StringVarP(&enumCmdInputPath, "input-path", "i", "", "eden generate enum --input-path=/go/src/eden-server")
	enumCmd.Flags().StringVarP(&enumCmdOutputPath, "output-path", "o", "", "eden generate enum --output-path=/go/src/eden-server")
	enumCmd.Flags().StringVarP(&enumCmdTypeName, "type-name", "t", "", "eden generate enum --type-name=Status")
}
