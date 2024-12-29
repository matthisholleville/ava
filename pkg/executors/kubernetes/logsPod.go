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

package kubernetes

import (
	"encoding/json"

	"github.com/matthisholleville/ava/pkg/common"
	v1 "k8s.io/api/core/v1"
)

var (
	tailLines = int64(100)
)

type LogsPod struct {
	PodName       string `json:"podName"`
	NamespaceName string `json:"namespaceName"`
}

func (LogsPod) Exec(e common.Executor, jsonString string) string {
	var podInfo LogsPod
	err := json.Unmarshal([]byte(jsonString), &podInfo)
	if err != nil {
		return "Error while retrieving the podName parameter:" + err.Error()
	}
	podLogOptions := v1.PodLogOptions{
		TailLines: &tailLines,
	}
	podLogs, err := e.Client.Client.CoreV1().Pods(podInfo.NamespaceName).GetLogs(podInfo.PodName, &podLogOptions).DoRaw(e.Context)
	if err != nil {
		return "Unable to retrieve pod logs." + err.Error()
	}

	return string(podLogs)

}
