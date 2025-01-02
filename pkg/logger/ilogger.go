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

package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/matthisholleville/ava/pkg/logger/json"
	"github.com/matthisholleville/ava/pkg/logger/raw"
	"go.uber.org/zap"
)

type LoggerWriter struct {
	Logger func(msg string, fields ...zap.Field)
}

func (lw *LoggerWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if line != "" {
			lw.Logger(line)
		}
	}
	return len(p), nil
}

type ILogger interface {
	Init(logLevel string) error
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func NewLogger(logType string) (ILogger, error) {
	switch logType {
	case "json":
		return &json.Json{}, nil
	case "raw":
		return &raw.Raw{}, nil
	default:
		return nil, fmt.Errorf("logger %s not found", logType)
	}
}

func InitLogger(logFormat, logLevel string) ILogger {
	logger, err := NewLogger(logFormat)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = logger.Init(logLevel)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return logger
}
