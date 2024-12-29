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
	"fmt"
	"os"

	"github.com/matthisholleville/ava/pkg/ai"
	"github.com/matthisholleville/ava/pkg/ai/openai"
	"github.com/matthisholleville/ava/pkg/logger"
)

// IKnowledge is the interface that wraps the basic methods to interact with a knowledge.
type IKnowledgeBackend interface {
	ConfigureKnowledge(logger logger.ILogger) error
	Purge() error
	GetName() string
	UploadFiles(path []string) error
}

type KnowledgeConfiguration struct {
	ActiveProvider string               `json:"activeProvider"`
	OpenAI         openai.Configuration `json:"openai"`
}

func NewBackendKnowledge(configuration KnowledgeConfiguration) (IKnowledgeBackend, error) {
	var knowledge IKnowledgeBackend

	switch configuration.ActiveProvider {
	case "openai":
		ai, err := ai.NewAI("openai", ai.WithPassword(os.Getenv("OPENAI_API_KEY")))
		if err != nil {
			return nil, err
		}
		knowledge = ai
	default:
		return nil, fmt.Errorf("provider %s not found", configuration.ActiveProvider)
	}
	return knowledge, nil
}
