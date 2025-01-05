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
	"time"

	"github.com/matthisholleville/ava/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RolloutDeployment struct {
	DeploymentName string `json:"deploymentName"`
	NamespaceName  string `json:"namespaceName"`
}

func (RolloutDeployment) GetName() string {
	return "rolloutDeployment"
}

func (RolloutDeployment) GetDescription() string {
	return "Perform a rollout restart for a deployment"
}

func (RolloutDeployment) GetParams() string {
	return `
	{
		"type": "object",
		"properties": {
			"deploymentName": {
				"type": "string"
			},
			"namespaceName": {
				"type": "string"
			}
		}
	}
	`
}

func (RolloutDeployment) Exec(e common.Executor, jsonString string) string {
	var rolloutInfo RolloutDeployment
	err := json.Unmarshal([]byte(jsonString), &rolloutInfo)
	if err != nil {
		return "Error while retrieving parameters: " + err.Error()
	}

	// Get the deployment
	client := e.Client.GetClient()
	deployment, err := client.AppsV1().Deployments(rolloutInfo.NamespaceName).Get(e.Context, rolloutInfo.DeploymentName, metav1.GetOptions{})
	if err != nil {
		return "Unable to retrieve deployment: " + err.Error()
	}

	// Modify the deployment to add the rollout restart annotation
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = map[string]string{}
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// Update the deployment
	_, err = client.AppsV1().Deployments(rolloutInfo.NamespaceName).Update(e.Context, deployment, metav1.UpdateOptions{})
	if err != nil {
		return "Failed to perform rollout restart: " + err.Error()
	}

	return "Rollout restart successfully triggered for deployment " + rolloutInfo.DeploymentName
}
