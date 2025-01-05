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

package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/matthisholleville/ava/pkg/common"
)

type GetUrl struct {
	Url string `json:"url"`
}

func (GetUrl) GetName() string {
	return "getUrl"
}

func (GetUrl) GetDescription() string {
	return "Call a URL and return the status code and response time"
}

func (GetUrl) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"url": {
			"type": "string"
			}
		}
	}
	`
}

func (GetUrl) Exec(e common.Executor, jsonString string) string {
	// execute a check on the url and return the status code
	var urlInfo GetUrl
	err := json.Unmarshal([]byte(jsonString), &urlInfo)
	if err != nil {
		return "Error while retrieving the url parameter:" + err.Error()
	}
	start := time.Now()

	resp, err := http.Get(urlInfo.Url)
	if err != nil {
		return "Error while executing the GET request:" + err.Error()
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return "GET request to " + urlInfo.Url + " returned status code " + resp.Status + " in " + duration.String()
}
