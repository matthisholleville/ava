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

type ListNamespaces struct{}

func (ListNamespaces) GetName() string {
	return "listNamespaces"
}

func (ListNamespaces) GetDescription() string {
	return "List all namespaces in the cluster"
}

func (ListNamespaces) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {}
	}
	`
}

func (ListNamespaces) Exec(e common.Executor, jsonString string) string {
	namespaces, err := e.Client.GetClient().CoreV1().Namespaces().List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list namespaces: " + err.Error()
	}
	result, _ := json.Marshal(namespaces)
	return string(result)
}
