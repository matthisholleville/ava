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

package json

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Json struct {
	*zap.Logger
}

func (j *Json) Info(msg string, fields ...zap.Field) {
	j.Logger.Info(msg, fields...)
}

func (j *Json) Debug(msg string, fields ...zap.Field) {
	j.Logger.Debug(msg, fields...)
}

func (j *Json) Warn(msg string, fields ...zap.Field) {
	j.Logger.Warn(msg, fields...)
}

func (j *Json) Error(msg string, fields ...zap.Field) {
	j.Logger.Error(msg, fields...)
}

func (j *Json) Fatal(msg string, fields ...zap.Field) {
	j.Logger.Fatal(msg, fields...)
}

func (j *Json) Init(logLevel string) error {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapEncoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(zapEncoderConfig, os.Stdout, level)
	j.Logger = zap.New(core)

	return nil
}
