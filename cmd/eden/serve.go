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
	"github.com/eden-framework/eden-framework/pkg/context"
	"github.com/eden-framework/eden-framework/pkg/courier"
	"github.com/eden-framework/eden-framework/pkg/courier/swagger"
	"github.com/eden-framework/eden-framework/pkg/courier/transport_http"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var serveCmdPort int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a OpenAPI server for your openapi.json",
	Run: func(cmd *cobra.Command, args []string) {

		var RootRouter = courier.NewRouter()
		RootRouter.Register(swagger.SwaggerRouter)

		server := &transport_http.ServeHTTP{
			Port:     serveCmdPort,
			WithCORS: true,
		}
		server.MarshalDefaults(server)

		ctx := context.NewWaitStopContext()

		go server.Serve(ctx, RootRouter)

		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)

		<-sig
		ctx.Cancel()
	},
}

func init() {
	serveCmd.Flags().IntVarP(&serveCmdPort, "port", "p", 8081, "eden serve --port=8081")
	rootCmd.AddCommand(serveCmd)
}
