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

package raw

import (
	"os"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

type Raw struct {
}

func (r *Raw) Info(msg string, fields ...zap.Field) {
	color.Green(msg)
}

func (r *Raw) Debug(msg string, fields ...zap.Field) {
	color.Cyan(msg)
}

func (r *Raw) Warn(msg string, fields ...zap.Field) {
	color.Yellow(msg)
}

func (r *Raw) Error(msg string, fields ...zap.Field) {
	color.Red(msg)
}

func (r *Raw) Fatal(msg string, fields ...zap.Field) {
	color.Red(msg)
	os.Exit(1)
}

func (r *Raw) Init(logLevel string) error {
	return nil
}
