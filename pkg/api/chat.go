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
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/ava/pkg/chat"
	"github.com/matthisholleville/ava/pkg/metrics"
	"go.uber.org/zap"
)

// CreateNewChat  Create program.
type CreateNewChat struct {
	Message  string `json:"message" example:"Pod web-server-5b866987d8-sxmtj in namespace default Crashlooping."`
	Language string `json:"language,omitempty" example:"en"`
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type WebhookPayload struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

// Chat godoc
// @Summary Chat with Ava from AlertManager Webhook
// @Description used to chat with Ava from AlertManager Webhook
// @Tags Chat
// @Accept json
// @Produce json
// @Router /chat/webhook [post]
//
//	@Param		_			body	WebhookPayload	true	"Webhook payload"
//
// @Success 202 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
func (s *Server) alertManagerWebhookChatHandler(echo echo.Context) error {
	s.logger.Info("Chatting with Ava")
	var data WebhookPayload

	if err := echo.Bind(&data); err != nil {
		s.logger.Error("reading the request body failed", zap.Error(err))
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusBadRequest)
	}

	// Check if there are alerts to process
	if !s.isFiring(data.Alerts) {
		s.logger.Info("No alerts to process")
		return s.JSONResponseWithCode(echo, "no alerts to process", http.StatusCreated)
	}

	chat, err := chat.NewChat(
		s.aiBackend,
		s.aiBackendPassword,
		s.logger,
		chat.WithLanguage("en"),
		chat.WithDbClient(s.db),
		chat.WithPersist(true),
		chat.WithConfigureAssistant(s.logger, s.enableExecutors),
	)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)

	}

	s.logger.Debug(fmt.Sprintf("Processing %d alerts", len(data.Alerts)))
	alreadyProcessedAlerts := []string{}
	for _, alert := range data.Alerts {
		if alert.Status != "firing" {
			continue
		}
		message := fmt.Sprintf("Summary: %s\nDescription: %s", alert.Annotations["summary"], alert.Annotations["description"])

		if slices.Contains(alreadyProcessedAlerts, message) {
			s.logger.Info(fmt.Sprintf("Alert already processed: %s", message))
			continue
		}

		alreadyProcessedAlerts = append(alreadyProcessedAlerts, message)
		s.logger.Info(fmt.Sprintf("Processing alert: %s", message))
		s.logger.Debug("Init Chat")
		threadID, err := chat.InitChat()
		if err != nil {
			s.logger.Fatal(err.Error())
			return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
		}

		go func() {
			chatType := "webhook"
			response, err := chat.Chat(message, threadID)
			if err != nil {
				metrics.ChatCounter.WithLabelValues("error", chatType).Inc()
				s.logger.Error(err.Error())
				return
			}
			s.logger.Info("Chat response processed successfully")
			metrics.ChatCounter.WithLabelValues("success", chatType).Inc()
			if chat.Persist {
				s.logger.Debug("Persisting chat")
				_, err := chat.PersistChat(message, response, threadID)
				if err != nil {
					s.logger.Error(err.Error())
					return
				}
				s.logger.Info("Chat saved successfully")
			}
		}()
	}

	return s.JSONResponseWithCode(echo, "alerts processed", http.StatusCreated)
}

// Check if alert fired
func (s *Server) isFiring(alerts []Alert) bool {
	for _, alert := range alerts {
		if alert.Status == "firing" {
			return true
		}
	}
	return false
}

// Chat godoc
// @Summary Chat with Ava
// @Description used to chat with Ava
// @Tags Chat
// @Accept json
// @Produce json
// @Router /chat [post]
//
//	@Param		_			body	CreateNewChat	true	"Create a new chat"
//
// @Success 202 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
func (s *Server) createChatHandler(echo echo.Context) error {
	s.logger.Info("Chatting with Ava")
	var data CreateNewChat

	if err := echo.Bind(&data); err != nil {
		s.logger.Error("reading the request body failed", zap.Error(err))
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusBadRequest)
	}

	chat, err := chat.NewChat(
		s.aiBackend,
		s.aiBackendPassword,
		s.logger,
		chat.WithLanguage(data.Language),
		chat.WithDbClient(s.db),
		chat.WithPersist(true),
		chat.WithConfigureAssistant(s.logger, s.enableExecutors),
	)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)

	}

	s.logger.Debug("Init Chat")
	threadID, err := chat.InitChat()
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	go func() {
		chatType := "chat"
		response, err := chat.Chat(data.Message, threadID)
		if err != nil {
			metrics.ChatCounter.WithLabelValues("error", chatType).Inc()
			s.logger.Error(err.Error())
			return
		}
		s.logger.Info("Chat response processed successfully")
		metrics.ChatCounter.WithLabelValues("success", chatType).Inc()
		if chat.Persist {
			s.logger.Debug("Persisting chat")
			_, err := chat.PersistChat(data.Message, response, threadID)
			if err != nil {
				s.logger.Error(err.Error())
				return
			}
			s.logger.Info("Chat saved successfully")
		}
	}()

	return s.JSONResponseWithCode(echo, fmt.Sprintf("/chat/%s", threadID), http.StatusCreated)
}

type FetchMessagesResponse struct {
	Chat     string `json:"chat"`
	Response string `json:"response"`
}

// Chat godoc
// @Summary Chat with Ava
// @Description used to chat with Ava
// @Tags Chat
// @Accept json
// @Produce json
// @Router /chat/{id} [get]
//
//	@Param		id	path	string				true	"ID"
//
// @Success 200 {object} FetchMessagesResponse
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
func (s *Server) fetchChatHandler(echo echo.Context) error {
	id := echo.Param("id")
	chat, err := chat.NewChat(
		s.aiBackend,
		s.aiBackendPassword,
		s.logger,
		chat.WithDbClient(s.db),
	)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}
	messages, err := chat.FetchChatMessages(id)
	if err != nil {
		s.logger.Error(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	s.logger.Debug(fmt.Sprintf("Preparing response. Number of messages: %d", len(messages)))
	responses := make([]FetchMessagesResponse, len(messages)-1)
	for _, message := range messages {
		response := FetchMessagesResponse{
			Chat:     message.Input,
			Response: message.Response,
		}
		responses = append(responses, response)
	}

	echo.JSONPretty(http.StatusOK, responses, "")
	return nil
}

// Chat godoc
// @Summary Chat with Ava
// @Description used to respond to Ava
// @Tags Chat
// @Accept json
// @Produce json
// @Router /chat/{id} [post]
//
//	@Param		id	path	string				true	"ID"
//	@Param		_			body	CreateNewChat	true	"Create a new chat"
//
// @Success 202 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
func (s *Server) respondChatHandler(echo echo.Context) error {
	id := echo.Param("id")

	var data CreateNewChat

	if err := echo.Bind(&data); err != nil {
		s.logger.Error("reading the request body failed", zap.Error(err))
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusBadRequest)
	}

	chat, err := chat.NewChat(
		s.aiBackend,
		s.aiBackendPassword,
		s.logger,
		chat.WithLanguage("en"),
		chat.WithDbClient(s.db),
		chat.WithPersist(true),
		chat.WithConfigureAssistant(s.logger, s.enableExecutors),
	)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	dbThread, err := chat.GetThread(id)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	go func() {
		chatType := "response"
		response, err := chat.Chat(data.Message, dbThread.ID)
		if err != nil {
			metrics.ChatCounter.WithLabelValues("error", chatType).Inc()
			s.logger.Error(err.Error())
			return
		}
		s.logger.Info("Chat response processed successfully")
		metrics.ChatCounter.WithLabelValues("success", chatType).Inc()
		if chat.Persist {
			s.logger.Debug("Persisting chat")
			_, err := chat.PersistChat(data.Message, response, dbThread.ID)
			if err != nil {
				s.logger.Error(err.Error())
				return
			}
			s.logger.Info("Chat saved successfully")
		}
	}()

	return s.JSONResponseWithCode(echo, fmt.Sprintf("/chat/%s", dbThread.ID), http.StatusCreated)

}
