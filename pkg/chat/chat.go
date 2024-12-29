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

package chat

import (
	"context"
	"time"

	db "github.com/matthisholleville/ava/internal/prisma"
	"github.com/matthisholleville/ava/pkg/ai"
	"github.com/matthisholleville/ava/pkg/common"
	"github.com/matthisholleville/ava/pkg/kubernetes"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/viper"
)

const (
	DEFAULT_SQL_TIMEOUT = 5 * time.Second
	DEFAULT_MAX_RESULTS = 1000
	DEFAULT_LANGUAGE    = "en"
)

type Chat struct {
	Context   context.Context
	Language  string
	AIClient  ai.IAI
	K8SClient *kubernetes.Client
	logger    logger.ILogger
	db        *db.PrismaClient
	Persist   bool
}

type Option func(*Chat)

func WithLanguage(language string) Option {

	if language == "" {
		language = DEFAULT_LANGUAGE
	}

	return func(i *Chat) {
		i.Language = language
	}
}

func WithDbClient(db *db.PrismaClient) Option {
	return func(i *Chat) {
		i.db = db
	}
}

func WithPersist(persist bool) Option {
	return func(i *Chat) {
		i.Persist = persist
	}
}

func WithConfigureAssistant(logger logger.ILogger) Option {
	return func(i *Chat) {
		err := i.AIClient.ConfigureAssistant(logger)
		if err != nil {
			i.logger.Fatal(err.Error())
		}
	}
}

func NewChat(
	backend string,
	password string,
	logger logger.ILogger,
	opts ...Option,
) (*Chat, error) {
	aiClient, err := ai.NewAI(backend, ai.WithPassword(password))
	if err != nil {
		return nil, err
	}

	kubecontext := viper.GetString("kubecontext")
	kubeconfig := viper.GetString("kubeconfig")
	k8sClient, err := kubernetes.NewClient(kubecontext, kubeconfig)
	if err != nil {
		return nil, err
	}

	client := &Chat{
		Context:   context.Background(),
		Language:  DEFAULT_LANGUAGE,
		AIClient:  aiClient,
		K8SClient: k8sClient,
		logger:    logger,
		db:        nil,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func (c *Chat) InitChat() (string, int, error) {
	c.logger.Info("Creates a new thread")
	threadID, err := c.AIClient.CreateThread()
	if err != nil {
		return "", -1, err
	}

	if c.Persist {
		c.logger.Info("Persists the message")
		persistedThread, err := c.PersistThread(*threadID)
		if err != nil {
			return "", -1, err
		}
		return *threadID, persistedThread.ID, nil
	}
	return *threadID, -1, nil
}

func (c *Chat) Chat(message, threadID string) (string, error) {

	c.logger.Info("Analyzes the message")
	return c.AIClient.Analyze(
		message,
		c.Language,
		threadID,
		common.Executor{
			Client:  c.K8SClient,
			Context: c.Context,
		},
	)
}

func (c *Chat) FetchChatMessages(threadID int) ([]db.ChatModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_SQL_TIMEOUT)
	defer cancel()
	return c.db.Chat.FindMany(
		db.Chat.ThreadID.Equals(threadID),
	).Take(DEFAULT_MAX_RESULTS).Exec(ctx)
}

func (c *Chat) PersistThread(threadID string) (*db.ThreadModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_SQL_TIMEOUT)
	defer cancel()
	return c.db.Thread.CreateOne(
		db.Thread.ThreadID.Set(threadID),
	).Exec(ctx)
}

func (c *Chat) FindThreadUnique(
	threadID string,
) (*db.ThreadModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_SQL_TIMEOUT)
	defer cancel()
	return c.db.
		Thread.
		FindUnique(
			db.Thread.ThreadID.Equals(threadID),
		).Exec(ctx)
}

func (c *Chat) PersistChat(chat, response, threadID string) (*db.ChatModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_SQL_TIMEOUT)
	defer cancel()
	return c.db.Chat.CreateOne(
		db.Chat.Input.Set(chat),
		db.Chat.Thread.Link(
			db.Thread.ThreadID.Equals(threadID),
		),
		db.Chat.Response.Set(response),
	).Exec(ctx)
}
