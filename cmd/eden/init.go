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
	"github.com/profzone/eden-framework/internal/project"
	"github.com/spf13/cobra"
	"os"
	"path"
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

		answers := struct {
			Name            string
			Group           string
			Owner           string
			Desc            string
			Version         string
			ProgramLanguage string `survey:"project_language"`
			Workflow        string
		}{}

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
		}

		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		newProject := currentProject.
			WithName(answers.Name).
			WithGroup(answers.Group).
			WithOwner(answers.Owner).
			WithDesc(answers.Desc).
			WithLanguage(answers.ProgramLanguage)

		if answers.Version != "" {
			newProject = newProject.WithVersion(answers.Version)
		}

		if answers.Workflow != "" && answers.Workflow != "custom" {
			newProject = newProject.WithWorkflow(answers.Workflow)
		}
		if newProject.Scripts == nil {
			newProject.Scripts = map[string]project.Script{
				"build": []string{"go build -v -o ./build/$PROJECT_NAME ./cmd && eden generate openapi"},
				"test":  []string{"go test ./cmd"},
			}
		}
		newProject.WriteToFile("./", "project.yml")
		newProject.SetEnviron()
		newProject.Workflow.TryExtendsOrSetDefaults().ToDroneConfig(&newProject).WriteToFile()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
