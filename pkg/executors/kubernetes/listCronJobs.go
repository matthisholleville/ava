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

type ListCronJobs struct {
	NamespaceName string `json:"namespaceName"`
}

func (ListCronJobs) GetName() string {
	return "listCronJobs"
}

func (ListCronJobs) GetDescription() string {
	return "List all cronjobs in the cluster"
}

func (ListCronJobs) GetParams() string {
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

func (ListCronJobs) Exec(e common.Executor, jsonString string) string {
	var cronInfos ListCronJobs
	err := json.Unmarshal([]byte(jsonString), &cronInfos)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	cronjobs, err := e.Client.GetClient().BatchV1().CronJobs(cronInfos.NamespaceName).List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list cronjobs: " + err.Error()
	}
	result, _ := json.Marshal(cronjobs)
	return string(result)
}
