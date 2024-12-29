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

package local

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matthisholleville/ava/pkg/logger"
)

type Local struct {
	logger logger.ILogger
}

func (l *Local) Configure(logger logger.ILogger) error {
	l.logger = logger
	return nil
}

func (l *Local) GetFiles(dir string) ([]string, error) {
	var files []string

	l.logger.Debug(fmt.Sprintf("Downloading files from %s", dir))
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".md" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func (l *Local) GetName() string {
	return "local"
}
