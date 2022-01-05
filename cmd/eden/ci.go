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
	"os"
	"regexp"
)

// ciCmd represents the ci command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "ci/cd workflow",
	Long:  fmt.Sprintf("%s\nci/cd workflow", CommandHelpHeader),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if currentProject.Feature == "" {
			featureName := getFeatureName()
			if featureName != "" {
				currentProject.Feature = featureName
			}
		}

		branchName := os.Getenv(project.EnvKeyCIBranch)
		if branchName != "master" {
			currentProject.Version.Prefix = currentProject.Feature
		}
		currentProject.Version.Suffix = getSha()
		currentProject.Selector = fmt.Sprintf("deployment-%s-%s", currentProject.Group, currentProject.Name)
		currentProject.SetEnviron()
	},
}

func init() {
	rootCmd.AddCommand(ciCmd)
}

var reFeatureBranch = regexp.MustCompile("feature/([a-z0-9\\-]+)")

func getFeatureName() string {
	matched := reFeatureBranch.FindAllStringSubmatch(os.Getenv(project.EnvKeyCIBranch), -1)
	if len(matched) > 0 {
		return matched[0][1]
	}
	return ""
}

func getSha() string {
	ref := os.Getenv(project.EnvKeyCICommitSHA)
	if len(ref) > 8 {
		return ref[0:8]
	}
	return ""
}
