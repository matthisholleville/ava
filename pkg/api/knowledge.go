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

	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/ava/pkg/ai/openai"
	backendKnowledge "github.com/matthisholleville/ava/pkg/knowledge/backend"
	sourceKnowledge "github.com/matthisholleville/ava/pkg/knowledge/source"
	"github.com/matthisholleville/ava/pkg/knowledge/source/configuration"
	"go.uber.org/zap"
)

// CreateNewKnowledge  Create knowledge.
type CreateNewKnowledge struct {
	Source           string `json:"source" example:"local"`
	Path             string `json:"path" example:"./docs/runbooks"`
	GitRepositoryURL string `json:"gitRepositoryURL,omitempty" example:""`
	GitAuthToken     string `json:"gitAuthToken,omitempty" example:""`
	GitBranch        string `json:"gitBranch,omitempty" example:""`
}

// Knowledge godoc
// @Summary Add knowledge to Ava
// @Description used to add knowledge to Ava
// @Tags Knowledge
// @Accept json
// @Produce json
// @Router /knowledge [post]
//
//	@Param		_			body	CreateNewKnowledge	true	"CreateNewKnowledge payload"
//
// @Success 202 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
func (s *Server) addKnowledgeHandler(echo echo.Context) error {
	s.logger.Info("Adding a new documents to Ava's knowledge base")
	var data CreateNewKnowledge

	if err := echo.Bind(&data); err != nil {
		s.logger.Error("reading the request body failed", zap.Error(err))
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusBadRequest)
	}

	backendKnowledge, err := backendKnowledge.NewBackendKnowledge(backendKnowledge.KnowledgeConfiguration{
		ActiveProvider: s.aiBackend,
		OpenAI: openai.Configuration{
			APIKey: s.aiBackendPassword,
		},
	})
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	s.logger.Info("Configuring backend knowledge")
	err = backendKnowledge.ConfigureKnowledge(s.logger)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	sourceKnowledge, err := sourceKnowledge.NewSourceKnowledge(data.Source)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	s.logger.Info("Configuring source knowledge")
	err = sourceKnowledge.Configure(s.logger, configuration.KnowledgeSourceConfiguration{
		Directory:        data.Path,
		GitRepositoryURL: data.GitRepositoryURL,
		GitAuthToken:     data.GitAuthToken,
		GitBranch:        data.GitBranch,
	})
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	defer sourceKnowledge.CleanUp()

	s.logger.Info("Getting files")
	files, err := sourceKnowledge.GetFiles()
	if err != nil {
		s.logger.Fatal(err.Error())
	}

	s.logger.Info(fmt.Sprintf("Found %d files", len(files)))

	s.logger.Info("Uploading file")
	err = backendKnowledge.UploadFiles(files)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	return s.JSONResponseWithCode(echo, "knowledge added", http.StatusCreated)
}

// PurgeKnowledge  Purge knowledge.
type PurgeKnowledge struct {
}

// Knowledge godoc
// @Summary Purge Ava's knowledge base
// @Description used to purge Ava's knowledge base
// @Tags Knowledge
// @Accept json
// @Produce json
// @Router /knowledge [delete]
//
//	@Param		_			body	PurgeKnowledge	true	"PurgeKnowledge payload"
//
// @Success 202 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
func (s *Server) purgeKnowledgeHandler(echo echo.Context) error {
	s.logger.Info("Purging Ava's knowledge base")
	var data CreateNewKnowledge

	if err := echo.Bind(&data); err != nil {
		s.logger.Error("reading the request body failed", zap.Error(err))
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusBadRequest)
	}

	backendKnowledge, err := backendKnowledge.NewBackendKnowledge(backendKnowledge.KnowledgeConfiguration{
		ActiveProvider: s.aiBackend,
		OpenAI: openai.Configuration{
			APIKey: s.aiBackendPassword,
		},
	})
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	s.logger.Info("Configuring backend knowledge")
	err = backendKnowledge.ConfigureKnowledge(s.logger)
	if err != nil {
		s.logger.Fatal(err.Error())
		return s.ErrorResponseWithCode(echo, err.Error(), http.StatusInternalServerError)
	}

	s.logger.Info("Purging knowledge")
	backendKnowledge.Purge()

	return s.JSONResponseWithCode(echo, "knowledge removed", http.StatusCreated)
}
