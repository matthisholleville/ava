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

package events

import (
	"fmt"

	db "github.com/matthisholleville/ava/internal/prisma"
	"github.com/matthisholleville/ava/pkg/events/slack"
	"github.com/matthisholleville/ava/pkg/logger"
)

var (
	clients = map[string]IEvent{
		"slack": &slack.SlackClient{},
	}
)

func GetClient(provider string) (IEvent, error) {
	if provider == "" {
		provider = "slack"
	}
	client, ok := clients[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", provider)
	}
	return client, nil
}

type IEvent interface {
	Configure(logger logger.ILogger, password string, db *db.PrismaClient) error
	SendMessage(channelID, message, ts string) error
	GetBotName(botID, teamID string) (string, error)
	ProcessEvent(data interface{}) (message string, threadID string, err error)
	PersistEvent(eventID, threadID string) (*db.EventModel, error)
	SendTechnicalErrorMessage(channelID, ts string) error
}
