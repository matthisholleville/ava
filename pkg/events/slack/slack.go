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

package slack

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	db "github.com/matthisholleville/ava/internal/prisma"
	"github.com/matthisholleville/ava/pkg/events/types"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

const (
	DEFAULT_SQL_TIMEOUT = 5 * time.Second
)

type SlackClient struct {
	*slack.Client
	logger logger.ILogger
	db     *db.PrismaClient
}

func (s *SlackClient) Configure(logger logger.ILogger, password string, db *db.PrismaClient) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	client := slack.New(password)
	s.Client = client
	s.logger = logger
	s.db = db
	return nil
}

func (s *SlackClient) SendMessage(channelID, message, ts string) error {
	_, _, err := s.Client.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionTS(ts),
	)
	return err
}

func (c *SlackClient) PersistEvent(eventID, threadID string) (*db.EventModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_SQL_TIMEOUT)
	defer cancel()
	return c.db.Event.CreateOne(
		db.Event.ID.Set(eventID),
		db.Event.Thread.Link(
			db.Thread.ID.Equals(threadID),
		),
	).Exec(ctx)
}

func (c *SlackClient) findUniqueEvent(
	eventID string,
) (*db.EventModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_SQL_TIMEOUT)
	defer cancel()
	return c.db.Event.FindUnique(
		db.Event.ID.Equals(eventID),
	).Exec(ctx)
}

func (s *SlackClient) extractUserName(message string) (username string, err error) {
	pattern := `^<@(.*?)>`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(message)

	if len(match) < 1 {
		return username, fmt.Errorf("no user found in message")
	}

	userID := match[1]

	userDetail, err := s.Client.GetUserInfo(userID)
	if err != nil {
		return username, err
	}

	return userDetail.Name, nil
}

func (s *SlackClient) processThread(eventData ReceiveSlackEvent) (message, threadID string) {

	// If the event is not a message, we ignore it
	if eventData.Event.ThreadTS == "" {
		return message, threadID
	}

	// If the event is a message, we extract the user name
	user, err := s.extractUserName(eventData.Event.Text)
	if err != nil {
		s.logger.Warn("Error extracting user name", zap.Error(err))
		return message, threadID
	}
	isMessageMentionAva := strings.Contains(user, "ava")

	// If the message does not mention Ava, we ignore it
	if !isMessageMentionAva {
		s.logger.Debug("Ignoring message not mentioning Ava")
		return message, threadID
	}

	// If the message mentions Ava, we find the thread
	event, err := s.findUniqueEvent(eventData.Event.ThreadTS)
	if err != nil {
		return message, threadID
	}
	threadID = event.ThreadID
	message = eventData.Event.Text

	s.logger.Debug(fmt.Sprintf("Found thread %s message with %s text", threadID, message))

	return message, threadID
}

func (s *SlackClient) ProcessEvent(data interface{}) (message string, threadID string, err error) {
	eventData := data.(ReceiveSlackEvent)
	var response string
	s.logger.Debug(fmt.Sprintf("Received event %s", eventData.Event.Subtype))

	isBotEvent := eventData.Event.BotID != ""

	// If the event is not a bot event, we check if it is a thread and if it mentions Ava
	if !isBotEvent {
		response, threadID = s.processThread(eventData)
		if response != "" {
			return response, threadID, nil
		}

		s.logger.Debug(fmt.Sprintf("Ignoring event %s", eventData.Event.Subtype))
		return response, threadID, fmt.Errorf("event ignored")
	}

	// If the event is a bot event, we check if it is an Ava bot event and we ignore it if it is
	isAvaBot, err := s.isAvaBotEvent(eventData.Event.BotID, eventData.TeamID)
	if isAvaBot || err != nil {
		s.logger.Debug(fmt.Sprintf("Ignoring event %s", eventData.Type))
		return response, threadID, fmt.Errorf("event ignored")
	}

	// If the event is a bot event and resolved message, we ignore it
	if eventData.Event.Attachments != nil {
		title := eventData.Event.Attachments[0].Title
		if s.isResolvedMessage(title) {
			s.logger.Debug("Ignoring resolved message")
			return response, threadID, fmt.Errorf("resolved message")
		}
	}

	response, err = s.extractMessage(eventData)
	if err != nil {
		return response, threadID, err
	}

	s.logger.Debug(fmt.Sprintf("Received message %s", response))

	return response, threadID, nil
}

func (s *SlackClient) isResolvedMessage(message string) bool {
	s.logger.Debug(fmt.Sprintf("Checking if message %s is resolved", message))
	return strings.Contains(message, "RESOLVED")
}

func (s *SlackClient) extractMessage(eventData ReceiveSlackEvent) (string, error) {
	message := eventData.Event.Text
	if message == "" && eventData.Event.Subtype == "bot_message" {
		if eventData.Event.Attachments != nil {
			message = eventData.Event.Attachments[0].Text
		}
	}

	if message == "" {
		return "", fmt.Errorf("no message found in event")
	}

	return message, nil
}

func (s *SlackClient) isAvaBotEvent(id, teamID string) (bool, error) {
	botName, err := s.GetBotName(id, teamID)
	if err != nil {
		return false, err
	}

	return strings.Contains(botName, "Ava"), nil
}

func (s *SlackClient) GetBotName(botID, teamID string) (string, error) {
	bot, err := s.Client.GetBotInfo(slack.GetBotInfoParameters{
		Bot:    botID,
		TeamID: teamID,
	})
	if err != nil {
		return "", err
	}

	return bot.Name, nil
}

func (s *SlackClient) SendTechnicalErrorMessage(channelID, ts string) error {
	return s.SendMessage(channelID, types.TechnicalErrorMessage, ts)
}
