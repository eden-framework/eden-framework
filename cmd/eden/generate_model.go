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
	mainGenerator "github.com/eden-framework/eden-framework/pkg/generator"
	"github.com/eden-framework/eden-framework/pkg/sqlx/generator"
	"github.com/spf13/cobra"
)

var cmdGenModelFlagDatabase string
var cmdGenModelFlagTableName string
var cmdGenModelFlagWithTableName bool
var cmdGenModelFlagWithTableInterfaces bool
var cmdGenModelFlagWithTableMethods bool
var cmdGenModelFlagWithComments bool

// apiCmd represents the api command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "generate crud operation for your model",
	Long:  fmt.Sprintf("%s\ngenerate model doc", CommandHelpHeader),
	Run: func(cmd *cobra.Command, args []string) {
		if cmdGenModelFlagDatabase == "" {
			panic("database must be defined")
		}

		for _, arg := range args {

			g := generator.NewSqlFuncGenerator()
			g.WithComments = true
			g.WithTableInterfaces = true
			g.StructName = arg
			g.Database = cmdGenModelFlagDatabase
			g.TableName = cmdGenModelFlagTableName
			g.WithTableName = cmdGenModelFlagWithTableName
			g.WithTableInterfaces = cmdGenModelFlagWithTableInterfaces
			g.WithMethods = cmdGenModelFlagWithTableMethods
			g.WithComments = cmdGenModelFlagWithComments

			mainGenerator.Generate(g, enumCmdInputPath, enumCmdOutputPath)
		}
	},
}

func init() {
	modelCmd.Flags().
		StringVarP(&cmdGenModelFlagDatabase, "database", "", "", "(required) register model to database var")
	modelCmd.Flags().
		StringVarP(&cmdGenModelFlagTableName, "table-name", "t", "", "custom table name")
	modelCmd.Flags().
		BoolVarP(&cmdGenModelFlagWithTableName, "with-table-name", "", true, "with Register and interface TableName")
	modelCmd.Flags().
		BoolVarP(&cmdGenModelFlagWithTableInterfaces, "with-table-interfaces", "", true, "with table interfaces like Indexes Fields")
	modelCmd.Flags().
		BoolVarP(&cmdGenModelFlagWithTableMethods, "with-methods", "", true, "with table methods")
	modelCmd.Flags().
		BoolVarP(&cmdGenModelFlagWithComments, "with-comments", "", false, "use comments")

	generateCmd.AddCommand(modelCmd)
}
