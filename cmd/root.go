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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/matthisholleville/ava/cmd/chat"
	"github.com/matthisholleville/ava/cmd/config"
	"github.com/matthisholleville/ava/cmd/knowledge"
	"github.com/matthisholleville/ava/cmd/serve"
	"github.com/matthisholleville/ava/internal/configuration"
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
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.PersistentFlags().StringVarP(&LogFormat, "log-format", "f", "raw", "Log format")
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "debug", "Log level")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("Default config file (%s/ava/ava.yaml)", xdg.ConfigHome))
}

func initConfig() {
	logger := logger.InitLogger(LogFormat, LogLevel)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println(xdg.ConfigHome)
		configDir := filepath.Join(xdg.ConfigHome, "ava")

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("ava")

		os.MkdirAll(configDir, 0755)

		configuration.WriteInitConfig(logger)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("unable to read ava configuration : %s.", cfgFile))
	}
	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	viper.Set("logger", logger)
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
