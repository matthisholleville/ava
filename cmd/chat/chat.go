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

package chat

import (
	"os"

	"github.com/matthisholleville/ava/pkg/chat"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	language    string
	backend     string
	kubecontext string
	kubeconfig  string
	message     string
)

var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with Ava",
	Long:  `Chat with Ava.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := viper.Get("logger").(logger.ILogger)

		logger.Info("Chatting with Ava")

		chat, err := chat.NewChat(
			backend,
			os.Getenv("OPENAI_API_KEY"),
			logger,
			chat.WithLanguage(language),
			chat.WithConfigureAssistant(logger),
		)
		if err != nil {
			logger.Fatal(err.Error())
		}

		threadID, _, err := chat.InitChat()
		if err != nil {
			logger.Fatal(err.Error())
		}

		response, err := chat.Chat(message, threadID)
		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info(response)

	},
}

func init() {
	ChatCmd.Flags().StringVarP(&message, "message", "m", "", "Message to send")
	ChatCmd.Flags().StringVarP(&language, "language", "g", "en", "Language to use")
	ChatCmd.Flags().StringVarP(&backend, "backend", "b", "openai", "Backend AI provider")
	ChatCmd.Flags().StringVar(&kubecontext, "kubecontext", "", "Kubernetes context to use. Only required if out-of-cluster.")
	ChatCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
}
