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

package api

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/ava/pkg/chat"
	"github.com/matthisholleville/ava/pkg/events/slack"
	"github.com/matthisholleville/ava/pkg/metrics"
	"go.uber.org/zap"
)

// Event godoc
// @Summary Receive a Slack event and chat with Ava
// @Description used to chat with Ava when a slack event is received
// @Tags Chat
// @Accept json
// @Produce json
// @Router /event/slack [post]
//
//	@Param		_			body	slack.ReceiveSlackEvent	true	"ReceiveSlackEvent payload"
//
// @Success 202 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
func (s *Server) slackEventHandler(echo echo.Context) error {
	s.logger.Info("Receiving a Slack event and chatting with Ava")

	var data slack.ReceiveSlackEvent
	if err := echo.Bind(&data); err != nil {
		s.logger.Error("reading the request body failed", zap.Error(err))
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusBadRequest)
	}

	slackValidationToken := os.Getenv("SLACK_VALIDATION_TOKEN")
	if slackValidationToken == "" || data.Token != slackValidationToken {
		s.logger.Error("Invalid slack token")
		return s.ErrorResponseWithCode(echo, "Invalid token", http.StatusBadRequest)
	}

	go func() {

		s.logger.Info("Processing event")
		var threadID string
		message, threadID, err := s.eventClient.ProcessEvent(data)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}

		// Check if executors are enabled and set default value to false
		isExecutorsEnabled := os.Getenv("ENABLE_EXECUTORS_ON_WEBHOOK") == "true"
		chat, err := chat.NewChat(
			"openai",
			os.Getenv("OPENAI_API_KEY"),
			s.logger,
			chat.WithLanguage("en"),
			chat.WithDbClient(s.db),
			chat.WithPersist(true),
			chat.WithConfigureAssistant(s.logger, isExecutorsEnabled),
		)
		if err != nil {
			s.logger.Fatal(err.Error())
			return
		}

		s.eventClient.SendMessage(data.Event.Channel, ":eyes:", data.Event.TS)

		// If threadID is empty, we need to initialize the chat
		if threadID == "" {
			s.logger.Debug("Init Chat")
			threadID, err = chat.InitChat()
			if err != nil {
				s.logger.Error("Init Chat failed", zap.Error(err))
				s.eventClient.SendTechnicalErrorMessage(data.Event.Channel, data.Event.TS)
				return
			}

			s.logger.Info("Persisting event")
			_, err := s.eventClient.PersistEvent(data.Event.EventTS, threadID)
			if err != nil {
				s.logger.Warn("Error persisting event", zap.Error(err))
			}
		}

		chatType := "slackEvent"
		response, err := chat.Chat(message, threadID)
		if err != nil {
			metrics.ChatCounter.WithLabelValues("error", chatType).Inc()
			s.logger.Error("Chat response processing failed", zap.Error(err))
			s.eventClient.SendTechnicalErrorMessage(data.Event.Channel, data.Event.TS)
			return
		}
		s.logger.Info("Chat response processed successfully")
		metrics.ChatCounter.WithLabelValues("success", chatType).Inc()
		if chat.Persist {
			s.logger.Debug("Persisting chat")
			_, err := chat.PersistChat(message, response, threadID)
			if err != nil {
				s.logger.Error("Chat saved failed", zap.Error(err))
				s.eventClient.SendTechnicalErrorMessage(data.Event.Channel, data.Event.TS)
				return
			}
			s.logger.Info("Chat saved successfully")
		}

		// Send the response to the slack channel
		err = s.eventClient.SendMessage(data.Event.Channel, response, data.Event.TS)
		if err != nil {
			s.logger.Error("sending message to slack failed", zap.Error(err))
			s.eventClient.SendTechnicalErrorMessage(data.Event.Channel, data.Event.TS)
			return
		}

	}()

	return echo.JSONPretty(http.StatusOK, data.Challenge, "")
}
