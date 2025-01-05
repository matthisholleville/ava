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
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	_ "github.com/matthisholleville/ava/docs"
	"github.com/matthisholleville/ava/internal/configuration"
	db "github.com/matthisholleville/ava/internal/prisma"
	"github.com/prometheus/client_golang/prometheus"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/ava/pkg/events"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/matthisholleville/ava/pkg/metrics"
	"go.uber.org/zap"
)

// @title Ava API
// @version 1.0
// @description This is the API for Ava, the SRE chatbot.

// @contact.name Source Code
// @contact.url https://github.com/matthisholleville/ava

// @BasePath /
// @schemes http https

var (
	healthy int32
	ready   int32
)

type Config struct {
	ServerShutdownTimeout time.Duration `mapstructure:"server-shutdown-timeout"`
	HttpServerTimeout     time.Duration `mapstructure:"http-server-timeout"`
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	PortMetrics           int           `mapstructure:"port-metrics"`
	Unhealthy             bool          `mapstructure:"unhealthy"`
	Unready               bool          `mapstructure:"unready"`
}

type Server struct {
	router            *echo.Echo
	logger            logger.ILogger
	config            *Config
	ctx               context.Context
	db                *db.PrismaClient
	eventClient       events.IEvent
	avaCfg            *configuration.Configuration
	aiBackend         string
	aiBackendPassword string
	enableExecutors   bool
}

func NewServer(config *Config, logger logger.ILogger, avaCfg *configuration.Configuration) (*Server, error) {

	dbClient := db.NewClient()
	if err := dbClient.Prisma.Connect(); err != nil {
		return nil, err
	}

	eventClient, err := events.GetClient(avaCfg.Events.Type)
	if err != nil {
		return nil, err
	}

	if avaCfg.Events.Type == "slack" {
		err = eventClient.Configure(logger, avaCfg.Events.Slack.BotToken, dbClient)
		if err != nil {
			return nil, err
		}
	}

	srv := &Server{
		router:            echo.New(),
		logger:            logger,
		config:            config,
		ctx:               context.Background(),
		db:                dbClient,
		eventClient:       eventClient,
		avaCfg:            avaCfg,
		aiBackend:         avaCfg.AI.Type,
		aiBackendPassword: avaCfg.AI.OpenAI.APIKey,
		enableExecutors:   avaCfg.Executors.Enabled,
	}

	return srv, nil
}

func (s *Server) registerHandlers() {
	if s.avaCfg.API.Swagger.Enabled {
		s.logger.Debug("Swagger enabled")
		s.router.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	s.router.GET("/live", s.healthzHandler)
	s.router.GET("/readyz", s.readyzHandler)

	if s.avaCfg.API.Chat.Enabled {
		s.logger.Debug("Chat API enabled")
		chat := s.router.Group("/chat")
		chat.POST("", s.createChatHandler)
		chat.POST("/webhook", s.alertManagerWebhookChatHandler)
		chat.GET("/:id", s.fetchChatHandler)
		chat.POST("/:id", s.respondChatHandler)
	}

	if s.avaCfg.API.Events.Enabled {
		s.logger.Debug("Events API enabled")
		event := s.router.Group("/event")
		if s.avaCfg.Events.Type == "slack" {
			event.POST("/slack", s.slackEventHandler)
		}
	}

	if s.avaCfg.API.Knowledge.Enabled {
		s.logger.Debug("Knowledge API enabled")
		knowledge := s.router.Group("/knowledge")
		knowledge.POST("", s.addKnowledgeHandler)
		knowledge.DELETE("", s.purgeKnowledgeHandler)
	}
}

func (s *Server) registerMiddlewares() {
	customMetrics := metrics.NewMetrics()
	customMetrics.RegisterCustomMetrics()
	s.router.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Registerer: prometheus.DefaultRegisterer,
		Namespace:  "ava",
	}))
}

func (s *Server) ListenAndServe() (*echo.Echo, *int32, *int32) {
	s.registerMiddlewares()
	s.registerHandlers()
	s.startMetricsServer()

	healthy = 1
	ready = 1

	return s.startServer(), &healthy, &ready

}

func (s *Server) startServer() *echo.Echo {

	s.router.HideBanner = true
	s.router.HidePort = true

	s.router.Server.WriteTimeout = s.config.HttpServerTimeout * time.Second
	s.router.Server.ReadTimeout = s.config.HttpServerTimeout * time.Second

	go func() {
		s.logger.Info("Starting HTTP Server.", zap.Int("port", s.config.Port))
		if err := s.router.Start(fmt.Sprintf(":%d", s.config.Port)); err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	if !s.config.Unhealthy {
		atomic.StoreInt32(&healthy, 1)
	}
	if !s.config.Unready {
		atomic.StoreInt32(&ready, 1)
	}

	return s.router

}

func (s *Server) startMetricsServer() {

	if s.config.PortMetrics > 0 {
		go func() {
			metrics := echo.New()
			metrics.HideBanner = true
			metrics.HidePort = true
			s.logger.Info("Starting Metrics HTTP Server.", zap.Int("port", s.config.PortMetrics))
			metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
			if err := metrics.Start(fmt.Sprintf(":%d", s.config.PortMetrics)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.logger.Fatal(err.Error())
			}
		}()
	}
}

type MapResponse map[string]string
