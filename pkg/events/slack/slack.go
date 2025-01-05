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
	"errors"
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

// reformatMessage reformat the message to remove the double asterisks
// that are not supported by Slack
func (s *SlackClient) reformatMessage(message string) string {
	message = strings.ReplaceAll(message, "#", "")
	return strings.ReplaceAll(message, "**", "*")
}

// SendMessage sends a message to a Slack channel
// with the given message and timestamp
// If the timestamp is empty, the message is sent as a new message
// If the timestamp is not empty, the message is sent as a reply to the message with the given timestamp
func (s *SlackClient) SendMessage(channelID, message, ts string) error {

	message = s.reformatMessage(message)

	params := slack.PostMessageParameters{
		Markdown: true,
	}
	_, _, err := s.Client.PostMessage(
		channelID,
		slack.MsgOptionText(message, true),
		slack.MsgOptionTS(ts),
		slack.MsgOptionPostMessageParameters(params),
	)
	return err
}

// PersistEvent persists an event in the database
// with the given eventID and threadID
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

// findUniqueEvent finds an event in the database
// with the given eventID
func (c *SlackClient) findUniqueEvent(
	eventID string,
) (*db.EventModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_SQL_TIMEOUT)
	defer cancel()
	return c.db.Event.FindUnique(
		db.Event.ID.Equals(eventID),
	).Exec(ctx)
}

// removeUserID removes the user ID from the message
// to keep only the message content
func (s *SlackClient) removeUserID(message string) string {
	pattern := `^<@.*?>\s*`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(message, "")
}

// extractUserName extracts the user name from the message
// by looking for the user ID
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

// processThread processes a thread event
// by looking for the message that mentions Ava
// and finding the thread ID
// If the message mentions Ava, the message content and thread ID are returned
func (s *SlackClient) processThread(eventData ReceiveSlackEvent) (message, threadID string) {

	// If the event is not a message, we ignore it
	if eventData.Event.ThreadTS == "" {
		return message, threadID
	}

	message = s.processAvaDirectMessage(eventData)
	if message == "" {
		return message, threadID
	}

	// If the message mentions Ava, we find the thread
	event, err := s.findUniqueEvent(eventData.Event.ThreadTS)
	if err != nil {
		return message, threadID
	}
	threadID = event.ThreadID

	s.logger.Debug(fmt.Sprintf("Found thread %s message with %s text", threadID, message))

	return message, threadID
}

// processAvaDirectMessage processes a direct message event
// by looking for the message that mentions Ava
// If the message mentions Ava, the message content is returned
func (s *SlackClient) processAvaDirectMessage(eventData ReceiveSlackEvent) (message string) {
	// If the event is a message, we extract the user name
	user, err := s.extractUserName(eventData.Event.Text)
	if err != nil {
		s.logger.Warn("Error extracting user name", zap.Error(err))
		return message
	}
	isMessageMentionAva := strings.Contains(user, "ava")

	// If the message does not mention Ava, we ignore it
	if !isMessageMentionAva {
		s.logger.Debug("Ignoring message not mentioning Ava")
		return message
	}

	return s.removeUserID(eventData.Event.Text)
}

func (s *SlackClient) ProcessEvent(data interface{}) (response string, threadID string, err error) {
	eventData := data.(ReceiveSlackEvent)
	s.logger.Debug(fmt.Sprintf("Received event %s", eventData.Event.Subtype))

	isBotEvent := eventData.Event.BotID != ""

	// If the event is not a bot event, we check if it is a thread and if it mentions Ava
	if !isBotEvent {

		response, threadID = s.processThread(eventData)
		if response != "" {
			return s.formatMessage(response), threadID, nil
		}

		response = s.processAvaDirectMessage(eventData)
		if response != "" {
			return response, threadID, nil
		}

		s.logger.Debug(fmt.Sprintf("Ignoring event %s", eventData.Event.Subtype))
		return response, threadID, fmt.Errorf("event ignored because it does not mention Ava or is not a thread managed by Ava")
	}

	// If the event is a bot event, we check if it is an Ava bot event and we ignore it if it is
	isAvaBot, err := s.isAvaBotEvent(eventData.Event.BotID, eventData.TeamID)
	if isAvaBot || err != nil {
		message := fmt.Sprintf("Ignoring event %s because the message is sent by ava", eventData.Event.Subtype)
		s.logger.Debug(message)
		return response, threadID, errors.New(message)
	}

	// If the event is a bot event and resolved message, we ignore it
	if eventData.Event.Attachments != nil {
		title := eventData.Event.Attachments[0].Title
		if s.isResolvedMessage(title) {
			message := fmt.Sprintf("Ignoring resolved message %s", title)
			s.logger.Debug(message)
			return response, threadID, errors.New(message)
		}
	}

	response, err = s.extractMessage(eventData)
	if err != nil {
		return response, threadID, err
	}

	s.logger.Debug(fmt.Sprintf("Received message %s", response))

	return s.formatMessage(response), threadID, nil
}

// isResolvedMessage checks if the message is a resolved message
// by looking for the word "RESOLVED" in the message
func (s *SlackClient) isResolvedMessage(message string) bool {
	s.logger.Debug(fmt.Sprintf("Checking if message %s is resolved", message))
	return strings.Contains(message, "RESOLVED")
}

// formatMessage formats the message to add a disclaimer
// that the output format must be Slack-friendly
func (s *SlackClient) formatMessage(message string) string {
	return fmt.Sprintf("%s.\n Do not use your default output format ! The output format of your response must be Slack-friendly. Only simple text, emojis, and Slack markdown format (*bold or title*, _italic_, `code inline`, and ```code blocks multi line```) are authorized. Your response must be clear and concise for the user.", message)

}

// extractMessage extracts the message from the event data
// by looking for the message in the event text
// If the message is not found, an error is returned
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

func (s *SlackClient) SendLookingMessage(channelID, ts string) error {
	return s.SendMessage(channelID, ":eyes:", ts)
}
