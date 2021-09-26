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
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/eden-framework/eden-framework/internal/generator"
	"github.com/eden-framework/eden-framework/internal/project"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var createCmdInitProject bool

// createCmd represents the create and init command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create and initialize a project",
	Run: func(cmd *cobra.Command, args []string) {
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

		// get the tag list of eden-framework
		//fmt.Println("fetch the tag list of eden-framework...")
		//cli := repo.NewClient("https", "api.github.com", 80)
		//tags, err := cli.GetTags("eden-framework/eden-framework")
		//if err != nil {
		//	logrus.Panicf("cannot get tag list of repo. err=%v", err)
		//}
		var tagList []string
		//for _, t := range tags {
		//	tagList = append(tagList, t.Name)
		//}
		tagList = append(tagList, "v1.1.9")
		if len(tagList) == 0 {
			logrus.Panic("cannot get tag list of repo. tag list is empty")
		}

		// get the plugin list of eden-framework
		//fmt.Println("fetch the plugin list of eden-framework...")
		//plugins, err := cli.GetPlugins()
		//if err != nil {
		//	logrus.Panicf("cannot get plugin list of eden-framework. err=%v", err)
		//}
		var pluginList []string
		var answers generator.ServiceOption
		//for _, p := range plugins {
		//	tags, err := cli.GetTags(p.FullName)
		//	if err != nil {
		//		logrus.Panicf("cannot get tag list of repo [%s]. err=%v", p.FullName, err)
		//	}
		//
		//	var pkgDisplayName, version string
		//	var pluginDetail generator.PluginDetail
		//	if len(tags) > 0 {
		//		version = tags[0].Name
		//		pkgDisplayName = fmt.Sprintf("%s@%s", p.GetPackagePath(), version)
		//		pluginDetail.Tag = tags[0]
		//	} else {
		//		pkgDisplayName = p.GetPackagePath()
		//	}
		//	pluginDetail.RepoFullName = p.FullName
		//	pluginDetail.PackageName = pkgDisplayName
		//	pluginDetail.PackagePath = p.GetPackagePath()
		//	pluginDetail.Version = version
		//	answers.PluginDetails = append(answers.PluginDetails, pluginDetail)
		//	pluginList = append(pluginList, pkgDisplayName)
		//}
		pluginList = append(pluginList, "abc")

		var qs = []*survey.Question{
			{
				Name: "framework_version",
				Prompt: &survey.Select{
					Message:  "框架版本",
					PageSize: 5,
					Options:  tagList,
					Default:  tagList[0],
				},
			},
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "项目名称",
					Default: currentProject.Name,
				},
				Validate: func(ans interface{}) error {
					err := survey.Required(ans)
					if err != nil {
						return err
					}

					cwd, _ := os.Getwd()
					p := path.Join(cwd, ans.(string))
					if generator.PathExist(p) {
						return fmt.Errorf("the path %s already exist", p)
					}
					return nil
				},
			},
			{
				Name: "package_name",
				Prompt: &survey.Input{
					Message: "包名",
				},
				Validate: survey.Required,
			},
			{
				Name: "database_support",
				Prompt: &survey.Select{
					Message: "数据库支持",
					Options: []string{"是", "否"},
					Default: "是",
				},
			},
			{
				Name: "apollo_support",
				Prompt: &survey.Select{
					Message: "Apollo配置中心支持",
					Options: []string{"是", "否"},
					Default: "是",
				},
			},
			{
				Name: "plugins",
				Prompt: &survey.MultiSelect{
					Message:  "插件支持",
					Options:  pluginList,
					PageSize: 5,
				},
			},
		}

		if createCmdInitProject {
			qs = append(qs, []*survey.Question{
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
						Default: currentProject.Workflow.Extends,
					},
				},
			}...)
		}

		err := survey.Ask(qs, &answers)
		if err != nil {
			if err == terminal.InterruptErr {
				return
			}
			panic(err)
		}

		cwd, _ := os.Getwd()

		gen := generator.NewServiceGenerator(answers)
		generator.Generate(gen, cwd, cwd)

		if createCmdInitProject {
			projectOpt := generator.ProjectOption{
				Name:            answers.Name,
				Group:           answers.Group,
				Owner:           answers.Owner,
				Desc:            answers.Desc,
				Version:         answers.Version,
				ProgramLanguage: answers.ProgramLanguage,
				Workflow:        answers.Workflow,
				ApolloSupport:   answers.ApolloSupport,
			}
			cwd = path.Join(cwd, answers.Name)
			projectGen := generator.NewProjectGenerator(projectOpt)
			generator.Generate(projectGen, cwd, cwd)
		}
	},
}

func init() {
	createCmd.Flags().BoolVarP(&createCmdInitProject, "init", "", true, "init after create")
	rootCmd.AddCommand(createCmd)
}
