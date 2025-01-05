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

type GetNode struct {
	NodeName string `json:"nodeName"`
}

func (GetNode) GetName() string {
	return "getNode"
}

func (GetNode) GetDescription() string {
	return "Get the details of a node"
}

func (GetNode) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"nodeName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetNode) Exec(e common.Executor, jsonString string) string {
	var nodeInfo GetNode
	err := json.Unmarshal([]byte(jsonString), &nodeInfo)
	if err != nil {
		return "Error while retrieving the nodeName parameter: " + err.Error()
	}
	node, err := e.Client.GetClient().CoreV1().Nodes().Get(e.Context, nodeInfo.NodeName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve node information: " + err.Error()
	}
	result, _ := json.Marshal(node)
	return string(result)
}
