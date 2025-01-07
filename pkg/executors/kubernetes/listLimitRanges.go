// Copyright © 2025 Ava AI.
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

package kubernetes

import (
	"encoding/json"

	"github.com/matthisholleville/ava/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListLimitRanges struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListLimitRanges) GetName() string {
	return "listLimitRanges"
}

func (ListLimitRanges) GetDescription() string {
	return "List all LimitRanges in a namespace"
}

func (ListLimitRanges) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (ListLimitRanges) Exec(e common.Executor, jsonString string) string {
	var lrInfo ListLimitRanges
	err := json.Unmarshal([]byte(jsonString), &lrInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	limitRanges, err := e.Client.GetClient().CoreV1().LimitRanges(lrInfo.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list LimitRanges: " + err.Error()
	}
	result, _ := json.Marshal(limitRanges)
	return string(result)
}
