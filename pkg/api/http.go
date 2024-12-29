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
	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSONResponseWithCode
func (s *Server) JSONResponseWithCode(echo echo.Context, result string, code int) error {
	data := SuccessResponse{
		Message: result,
	}
	return echo.JSONPretty(code, data, "")
}

// ErrorResponseWithCode
func (s *Server) ErrorResponseWithCode(echo echo.Context, error string, code int) error {
	data := ErrorResponse{
		Code:    code,
		Message: error,
	}

	return echo.JSONPretty(code, data, "")
}
