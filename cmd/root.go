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
	"os"
	"strings"

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
	cfgFile   string
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config/", "config file (default is $pwd/config/ava.yaml)")
}

func initConfig() {
	logger := logger.InitLogger(LogFormat, LogLevel)
	viper.Set("logger", logger)

	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfgFile)
	viper.SetConfigName("ava.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic("unable to read ava configuration.")
	}
	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
