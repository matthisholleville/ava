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
	"os/exec"

	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Edit the configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := viper.Get("logger").(logger.ILogger)
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			logger.Fatal("No configuration file found to dump.")
		}
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi" // Default to vi if EDITOR is not set
		}

		c := exec.Command(editor, configFile)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			logger.Fatal(fmt.Sprintf("Failed to edit configuration file: %v", err))
		}

		logger.Info("Configuration updated successfully.")
	},
}

func init() {
}
