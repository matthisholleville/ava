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

package knowledge

import (
	"os"

	"github.com/matthisholleville/ava/pkg/ai/openai"
	knowledge "github.com/matthisholleville/ava/pkg/knowledge/backend"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge Ava's knowledge base",
	Long:  `Purge Ava's knowledge base.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := viper.Get("logger").(logger.ILogger)
		logger.Info("Purging Ava's knowledge base")
		knowledge, err := knowledge.NewBackendKnowledge(knowledge.KnowledgeConfiguration{
			ActiveProvider: "openai",
			OpenAI: openai.Configuration{
				APIKey: os.Getenv("OPENAI_API_KEY"),
			},
		})
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Configuring knowledge")
		err = knowledge.ConfigureKnowledge(logger)
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Purging knowledge")
		knowledge.Purge()
	},
}
