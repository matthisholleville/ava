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

package ai

import (
	"fmt"

	"github.com/matthisholleville/ava/pkg/ai/openai"
	"github.com/matthisholleville/ava/pkg/common"
	"github.com/matthisholleville/ava/pkg/logger"
)

type Option func(*AIProvider)

func WithPassword(password string) Option {
	return func(a *AIProvider) {
		a.password = password
	}
}

type IAI interface {
	Configure(logger logger.ILogger) error
	ConfigureKnowledge(logger logger.ILogger) error
	ConfigureAssistant(logger logger.ILogger) error
	Purge() error
	UploadFiles(path []string) error
	GetName() string
	CreateThread() (*string, error)
	Analyze(text, language string, threadID string, executorConfig common.Executor) (string, error)
}

type AIProvider struct {
	password string
}

func NewAI(backend string, opts ...Option) (IAI, error) {
	var ai IAI

	aiProvider := &AIProvider{}

	for _, opt := range opts {
		opt(aiProvider)
	}
	switch backend {
	case "openai":
		ai = &openai.OpenAIClient{
			Configuration: openai.Configuration{
				APIKey: aiProvider.password,
			},
		}

	default:
		return nil, fmt.Errorf("backend %s not found", backend)
	}
	return ai, nil
}
