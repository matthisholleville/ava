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

package signals

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Shutdown struct {
	logger                logger.ILogger
	serverShutdownTimeout time.Duration
}

func NewShutdown(serverShutdownTimeout time.Duration, logger logger.ILogger) (*Shutdown, error) {
	srv := &Shutdown{
		logger:                logger,
		serverShutdownTimeout: serverShutdownTimeout,
	}

	return srv, nil
}

func (s *Shutdown) Graceful(stopCh <-chan struct{}, httpServer *echo.Echo, healthy *int32, ready *int32) {
	ctx := context.Background()

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(ctx, s.serverShutdownTimeout)
	defer cancel()

	// all calls to /healthz and /readyz will fail from now on
	atomic.StoreInt32(healthy, 0)
	atomic.StoreInt32(ready, 0)

	s.logger.Info("Shutting down HTTP server", zap.Duration("timeout", s.serverShutdownTimeout))

	// There could be a period where a terminating pod may still receive requests. Implementing a brief wait can mitigate this.
	// See: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination
	// the readiness check interval must be lower than the timeout
	if viper.GetString("level") != "debug" {
		time.Sleep(3 * time.Second)
	}

	// determine if the http server was started
	if httpServer != nil {
		if err := httpServer.Shutdown(ctx); err != nil {
			s.logger.Warn("HTTP server graceful shutdown failed", zap.Error(err))
		}
	}
}
