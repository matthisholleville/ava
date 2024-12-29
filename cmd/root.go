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

package cmd

import (
	"github.com/matthisholleville/ava/cmd/chat"
	"github.com/matthisholleville/ava/cmd/knowledge"
	"github.com/matthisholleville/ava/cmd/serve"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version   string
	Commit    string
	Date      string
	LogFormat string
	LogLevel  string
)

var rootCmd = &cobra.Command{
	Use:   "ava",
	Short: "Ava is a SRE AI assistant",
	Long:  `Ava is a SRE AI assistant that helps you to monitor your infrastructure and applications.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(knowledge.KnowledgeCmd)
	rootCmd.AddCommand(chat.ChatCmd)
	rootCmd.AddCommand(serve.ServeCmd)
	rootCmd.PersistentFlags().StringVarP(&LogFormat, "log-format", "f", "raw", "Log format")
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "debug", "Log level")
}

func initConfig() {
	logger := logger.InitLogger(LogFormat, LogLevel)
	viper.Set("logger", logger)
}
