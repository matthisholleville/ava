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

	"github.com/matthisholleville/ava/pkg/knowledge/source/configuration"
	"github.com/matthisholleville/ava/pkg/knowledge/source/git"
	"github.com/matthisholleville/ava/pkg/knowledge/source/local"
	"github.com/matthisholleville/ava/pkg/logger"
)

type IKnowledgeSource interface {
	Configure(logger logger.ILogger, config configuration.KnowledgeSourceConfiguration) error
	GetName() string
	GetFiles() ([]string, error)
	CleanUp() error
}

func NewSourceKnowledge(provider string) (IKnowledgeSource, error) {
	var knowledge IKnowledgeSource

	switch provider {
	case "local":
		knowledge = &local.Local{}
	case "git":
		knowledge = &git.Git{}
	default:
		return nil, fmt.Errorf("provider %s not found", provider)
	}
	return knowledge, nil
}
