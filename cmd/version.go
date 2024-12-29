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
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Ava",
	Long:  `All software has versions. This is Ava's`,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "dev" {
			details, ok := debug.ReadBuildInfo()
			if ok && details.Main.Version != "" && details.Main.Version != "(devel)" {
				Version = details.Main.Version
				for _, i := range details.Settings {
					if i.Key == "vcs.time" {
						Date = i.Value
					}
					if i.Key == "vcs.revision" {
						Commit = i.Value
					}
				}
			}
		}
		fmt.Printf("ava: %s (%s), built at: %s\n", Version, Commit, Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
