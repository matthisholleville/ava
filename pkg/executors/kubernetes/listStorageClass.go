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

type ListStorageClasses struct{}

func (ListStorageClasses) GetName() string {
	return "listStorageClasses"
}

func (ListStorageClasses) GetDescription() string {
	return "List all StorageClasses in the cluster"
}

func (ListStorageClasses) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {}
	}
	`
}

func (ListStorageClasses) Exec(e common.Executor, jsonString string) string {
	storageClasses, err := e.Client.GetClient().StorageV1().StorageClasses().List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list StorageClasses: " + err.Error()
	}
	result, _ := json.Marshal(storageClasses)
	return string(result)
}
