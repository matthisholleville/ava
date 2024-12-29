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
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

// Healthz godoc
// @Summary Liveness check
// @Description used by Kubernetes liveness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /live [get]
// @Success 200 {string} string "OK"
// @Success 503 {string} string "KO"
func (s *Server) healthzHandler(echo echo.Context) error {
	if atomic.LoadInt32(&healthy) == 1 {
		echo.JSONPretty(http.StatusOK, map[string]string{"status": "OK"}, "")
		return nil
	}
	echo.JSONPretty(http.StatusServiceUnavailable, map[string]string{"status": "KO"}, "")
	return nil
}

// Readyz godoc
// @Summary Readiness check
// @Description used by Kubernetes readiness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz [get]
// @Success 200 {string} string "OK"
// @Success 503 {string} string "KO"
func (s *Server) readyzHandler(echo echo.Context) error {
	if atomic.LoadInt32(&ready) == 1 {
		echo.JSONPretty(http.StatusOK, map[string]string{"status": "OK"}, "")
		return nil
	}
	echo.JSONPretty(http.StatusServiceUnavailable, map[string]string{"status": "KO"}, "")
	return nil
}
