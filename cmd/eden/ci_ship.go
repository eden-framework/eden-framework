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
	"github.com/eden-framework/eden-framework/internal/project"
	"github.com/spf13/cobra"
)

var shipFlagPush bool

// ciShipCmd represents the ciShip command
var ciShipCmd = &cobra.Command{
	Use:   "ship",
	Short: "ci ship a project as a image",
	Run: func(cmd *cobra.Command, args []string) {
		currentProject.Run(project.CommandsForShipping(currentProject, shipFlagPush)...)
	},
}

func init() {
	ciShipCmd.Flags().BoolVarP(&shipFlagPush, "push", "", false, "push after build")
	ciCmd.AddCommand(ciShipCmd)
}
