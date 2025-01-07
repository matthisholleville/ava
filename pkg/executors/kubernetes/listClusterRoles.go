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

type ListClusterRoles struct{}

func (ListClusterRoles) GetName() string {
	return "listClusterRoles"
}

func (ListClusterRoles) GetDescription() string {
	return "List all ClusterRoles in the cluster"
}

func (ListClusterRoles) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {}
	}
	`
}

func (ListClusterRoles) Exec(e common.Executor, jsonString string) string {
	clusterRoles, err := e.Client.GetClient().RbacV1().ClusterRoles().List(e.Context, metav1.ListOptions{})
	if err != nil {
		return "Unable to list ClusterRoles: " + err.Error()
	}
	result, _ := json.Marshal(clusterRoles)
	return string(result)
}
