// Copyright © 2025 Ava AI.
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

package config

import (
	"fmt"
	"os"

	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ()

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump the configuration file",
	Long:  `Dump the configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := viper.Get("logger").(logger.ILogger)
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			logger.Fatal("No configuration file found to dump.")
		}
		content, err := os.ReadFile(configFile)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Unable to read configuration file: %v", err))
		}

		logger.Info(string(content))
	},
}

func init() {
}
