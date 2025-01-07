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

type GetPersistentVolume struct {
	PersistentVolumeName string `json:"persistentVolumeName"`
}

func (GetPersistentVolume) GetName() string {
	return "getPersistentVolume"
}

func (GetPersistentVolume) GetDescription() string {
	return "Retrieve details of a PersistentVolume"
}

func (GetPersistentVolume) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"persistentVolumeName": {
				"type": "string"
			}
		}
	}
	`
}

func (GetPersistentVolume) Exec(e common.Executor, jsonString string) string {
	var pvInfo GetPersistentVolume
	err := json.Unmarshal([]byte(jsonString), &pvInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}
	persistentVolume, err := e.Client.GetClient().CoreV1().PersistentVolumes().Get(e.Context, pvInfo.PersistentVolumeName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve PersistentVolume information: " + err.Error()
	}
	result, _ := json.Marshal(persistentVolume)
	return string(result)
}
