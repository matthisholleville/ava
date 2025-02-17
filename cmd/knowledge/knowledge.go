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

import "github.com/spf13/cobra"

var (
	backend string
	path    string
	source  string
)

var KnowledgeCmd = &cobra.Command{
	Use:   "knowledge",
	Short: "Enrich Ava's knowledge base",
	Long:  `Enrich Ava's knowledge base with new information.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	KnowledgeCmd.AddCommand(addCmd)
	KnowledgeCmd.AddCommand(purgeCmd)
}
