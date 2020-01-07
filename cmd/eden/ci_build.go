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
	"github.com/spf13/cobra"
)

// ciBuildCmd represents the ciBuild command
var ciBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "ci build workflow",
	Long:  fmt.Sprintf("%s\ngenerate api doc", CommandHelpHeader),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	ciCmd.AddCommand(ciBuildCmd)
}
