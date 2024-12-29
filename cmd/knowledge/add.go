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
	backendKnowledge "github.com/matthisholleville/ava/pkg/knowledge/backend"
	sourceKnowledge "github.com/matthisholleville/ava/pkg/knowledge/source"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultBackend    = "openai"
	defaultSourcePath = "./docs/runbooks"
	defaultSource     = "local"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new document to Ava's knowledge base",
	Long:  `Add a new document to Ava's knowledge base.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := viper.Get("logger").(logger.ILogger)
		logger.Info("Adding a new document to Ava's knowledge base")
		backendKnowledge, err := backendKnowledge.NewBackendKnowledge(backendKnowledge.KnowledgeConfiguration{
			ActiveProvider: "openai",
			OpenAI: openai.Configuration{
				APIKey: os.Getenv("OPENAI_API_KEY"),
			},
		})
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Configuring backend knowledge")
		err = backendKnowledge.ConfigureKnowledge(logger)
		if err != nil {
			logger.Fatal(err.Error())
		}

		sourceKnowledge, err := sourceKnowledge.NewSourceKnowledge(source)
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Configuring source knowledge")
		sourceKnowledge.Configure(logger)

		logger.Info("Getting files")
		files, err := sourceKnowledge.GetFiles(path)
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Uploading file")
		err = backendKnowledge.UploadFiles(files)
		if err != nil {
			logger.Fatal(err.Error())
		}
	},
}

func init() {
	addCmd.Flags().StringVarP(&backend, "backend", "b", defaultBackend, "Backend AI provider")
	addCmd.Flags().StringVarP(&path, "path", "p", defaultSourcePath, "Path to the document to add")
	addCmd.Flags().StringVarP(&source, "source", "s", defaultSource, "Source of the document")
}
