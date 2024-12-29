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

package openai

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

var (
	k8sPodParam = `
	{
		"type": "object",
		"properties": {
			"podName": {
			"type": "string"
			},
			"namespaceName": {
			"type": "string"
			}
		}
	}
	`
	urlParam = `
	{
		"type": "object",
		"properties": {
			"url": {
			"type": "string"
			}
		}
	}
	`

	waitParam = `
	{
		"type": "object",
		"properties": {
			"time": {
			"type": "string"
			}
		}
	}
	`
	k8sFunctions = []openai.AssistantTool{
		{
			Function: &openai.FunctionDefinition{
				Name:        "podLogs",
				Parameters:  json.RawMessage([]byte(k8sPodParam)),
				Description: "Get the logs of a pod",
			},
			Type: openai.AssistantToolTypeFunction,
		},
		{
			Function: &openai.FunctionDefinition{
				Name:        "getPod",
				Parameters:  json.RawMessage([]byte(k8sPodParam)),
				Description: "Get the details of a pod",
			},
			Type: openai.AssistantToolTypeFunction,
		},
		{
			Function: &openai.FunctionDefinition{
				Name:        "deletePod",
				Parameters:  json.RawMessage([]byte(k8sPodParam)),
				Description: "Delete a pod",
			},
			Type: openai.AssistantToolTypeFunction,
		},
		{
			Function: &openai.FunctionDefinition{
				Name:        "getUrl",
				Parameters:  json.RawMessage([]byte(urlParam)),
				Description: "Call a URL and return the status code and response time",
			},
			Type: openai.AssistantToolTypeFunction,
		},
		{
			Function: &openai.FunctionDefinition{
				Name:        "wait",
				Parameters:  json.RawMessage([]byte(waitParam)),
				Description: "Wait for a specified time in seconds",
			},
			Type: openai.AssistantToolTypeFunction,
		},
	}

	Functions = k8sFunctions
)
