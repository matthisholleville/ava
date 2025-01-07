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

type GetStorageClass struct {
	StorageClassName string `json:"storageClassName"`
}

func (GetStorageClass) GetName() string {
	return "getStorageClass"
}

func (GetStorageClass) GetDescription() string {
	return "Retrieve details of a StorageClass"
}

func (GetStorageClass) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"storageClassName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetStorageClass) Exec(e common.Executor, jsonString string) string {
	var scInfo GetStorageClass
	err := json.Unmarshal([]byte(jsonString), &scInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	storageClass, err := e.Client.GetClient().StorageV1().StorageClasses().Get(e.Context, scInfo.StorageClassName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve StorageClass information: " + err.Error()
	}
	result, _ := json.Marshal(storageClass)
	return string(result)
}
