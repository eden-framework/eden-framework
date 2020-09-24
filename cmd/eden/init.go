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
	"github.com/AlecAivazis/survey/v2"
	"github.com/eden-framework/eden-framework/internal/project"
	"github.com/eden-framework/eden-framework/pkg/generator"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(currentProject.Name) == 0 {
			cwd, _ := os.Getwd()
			currentProject.Name = path.Base(cwd)
		}

		if len(currentProject.Group) == 0 {
			currentProject.Group = "profzone"
		}

		if len(currentProject.Owner) == 0 {
			currentProject.Owner = "profzone"
		}

		if len(currentProject.ProgramLanguage) == 0 {
			currentProject.ProgramLanguage = "golang"
		}

		if currentProject.Workflow.Extends == "" {
			currentProject.Workflow.Extends = "feature-pr"
		}

		answers := generator.ProjectOption{}

		var qs = []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "项目名称",
					Default: currentProject.Name,
				},
				Validate: survey.Required,
			},
			{
				Name: "desc",
				Prompt: &survey.Input{
					Message: "项目描述",
					Default: currentProject.Desc,
				},
				Validate: survey.Required,
			},
			{
				Name: "group",
				Prompt: &survey.Input{
					Message: "项目所属应用",
					Default: currentProject.Group,
				},
				Validate: survey.Required,
			},
			{
				Name: "owner",
				Prompt: &survey.Input{
					Message: "项目所属用户组",
					Default: currentProject.Owner,
				},
				Validate: survey.Required,
			},
			{
				Name: "version",
				Prompt: &survey.Input{
					Message: "项目版本号 (x.x.x)",
					Default: currentProject.Version.String(),
				},
				Validate: survey.Required,
			},
			{
				Name: "project_language",
				Prompt: &survey.Select{
					Message: "项目所用编程语言",
					Options: append(project.RegisteredBuilders.SupportProgramLanguages(), "custom"),
					Default: currentProject.ProgramLanguage,
				},
			},
			{
				Name: "workflow",
				Prompt: &survey.Select{
					Message: "项目 workflow",
					Options: append(project.PresetWorkflows.List(), "custom"),
					Default: func() string {
						if currentProject.Workflow.Extends == "" {
							return "feature-pr"
						}
						return currentProject.Workflow.Extends
					}(),
				},
			},
			{
				Name: "apollo_support",
				Prompt: &survey.Select{
					Message: "是否启用Apollo配置中心支持",
					Options: []string{"是", "否"},
					Default: func() string {
						if script, ok := currentProject.Scripts["build"]; ok {
							if strings.Contains(script.String(), "apollo.Branch") {
								return "是"
							}
						}
						return "否"
					}(),
				},
			},
		}

		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		cwd, _ := os.Getwd()

		gen := generator.NewProjectGenerator(answers)
		generator.Generate(gen, cwd, cwd)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
