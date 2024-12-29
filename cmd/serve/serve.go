// Copyright © 2024 Ava AI.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package serve

import (
	"time"

	"github.com/matthisholleville/ava/pkg/api"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/matthisholleville/ava/pkg/signals"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	httpServerTimeout     time.Duration
	port                  int
	metricsPort           int
	unhealthy             bool
	ready                 bool
	serverShutdownTimeout time.Duration
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve Ava API",
	Long:  `Serve Ava API.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := viper.Get("logger").(logger.ILogger)

		logger.Info("Serving Ava API")

		// Start the server
		serverConfig := &api.Config{
			HttpServerTimeout: httpServerTimeout,
			Port:              port,
			PortMetrics:       metricsPort,
			Host:              "",
			Unhealthy:         unhealthy,
			Unready:           ready,
		}

		server, _ := api.NewServer(serverConfig, logger)
		httpServer, healthy, ready := server.ListenAndServe()

		// graceful shutdown
		stopCh := signals.SetupSignalHandler()
		sd, _ := signals.NewShutdown(serverShutdownTimeout, logger)
		sd.Graceful(stopCh, httpServer, healthy, ready)
	},
}

func init() {
	ServeCmd.Flags().DurationVar(&httpServerTimeout, "http-server-timeout", 10*time.Second, "HTTP server timeout")
	ServeCmd.Flags().IntVar(&port, "port", 8080, "Port to listen on")
	ServeCmd.Flags().IntVar(&metricsPort, "port-metrics", 8081, "Port to listen on for metrics")
	ServeCmd.Flags().BoolVar(&unhealthy, "unhealthy", true, "Set the server as unhealthy")
	ServeCmd.Flags().BoolVar(&ready, "ready", true, "Set the server as ready")
	ServeCmd.Flags().DurationVar(&serverShutdownTimeout, "server-shutdown-timeout", 30*time.Second, "Server shutdown timeout")
}
