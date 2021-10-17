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
	"github.com/eden-framework/eden-framework/internal/project"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func main() {
	Execute()
}

const CommandHelpHeader = `

   _______  _____  __  ____                                   __
  / __/ _ \/ __/ |/ / / __/______ ___ _  ___ _    _____  ____/ /__
 / _// // / _//    / / _// __/ _ '/  ' \/ -_) |/|/ / _ \/ __/  '_/
/___/____/___/_/|_/ /_/ /_/  \_,_/_/_/_/\__/|__,__/\___/_/ /_/\_\


eden-framework staging tool chain
`

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eden",
	Short: "eden-framework staging tool chain",
	Long:  CommandHelpHeader,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eden.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".eden-framework" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".eden")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	initProject()
}

func initProject() {
	project.RegisterBuilder("BUILDER_RANCHER", &project.Builder{
		Image:      getConfigOrDefault("PROJECT_BUILDER_RANCHER", "profzone/golang-ondeploy:1.0"),
		WorkingDir: "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.RegisterBuilder("BUILDER_DOCKER", &project.Builder{
		Image:      getConfigOrDefault("PROJECT_BUILDER_DOCKER", "profzone/golang-onship:1.14"),
		WorkingDir: "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.RegisterBuilder("BUILDER_GOLANG", &project.Builder{
		ProgramLanguage: "golang",
		Image:           getConfigOrDefault("PROJECT_BUILDER_GOLANG", "profzone/golang-onbuild:1.14"),
		WorkingDir:      "/go/src/github.com/${PROJECT_GROUP}/${PROJECT_NAME}",
	})
	project.DockerRegistry = getConfigOrDefault("PROJECT_DOCKER_REGISTRY", "registry.cn-hangzhou.aliyuncs.com")
}
