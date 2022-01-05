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

var errCmdInputPath, errCmdOutputPath string

// errorCmd represents the enum command
var errorCmd = &cobra.Command{
	Use:   "error",
	Short: "generate error's go files",
	Long:  fmt.Sprintf("%s\ngenerate error", CommandHelpHeader),
	Run: func(cmd *cobra.Command, args []string) {
		if errCmdInputPath == "" {
			errCmdInputPath, _ = os.Getwd()
		}
		if errCmdOutputPath == "" {
			errCmdOutputPath, _ = os.Getwd()
		}
		gen := generator.NewStatusErrGenerator()

		generator.Generate(gen, errCmdInputPath, errCmdOutputPath)
	},
}

func init() {
	generateCmd.AddCommand(errorCmd)
	errorCmd.Flags().StringVarP(&errCmdInputPath, "input-path", "i", "", "eden generate error --input-path=/go/src/eden-server")
	errorCmd.Flags().StringVarP(&errCmdOutputPath, "output-path", "o", "", "eden generate error --output-path=/go/src/eden-server")
}
