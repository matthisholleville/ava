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

package executors

import (
	"github.com/matthisholleville/ava/internal/configuration"
	"github.com/matthisholleville/ava/pkg/common"
	commonExecutorsPkg "github.com/matthisholleville/ava/pkg/executors/common"
	"github.com/matthisholleville/ava/pkg/executors/kubernetes"
	"github.com/matthisholleville/ava/pkg/executors/web"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/spf13/viper"
)

var (
	k8sReadExecutors = map[string]IExecutor{
		"describeService": kubernetes.DescribeService{},
		"getCronJobs":     kubernetes.GetCronJobs{},
		"getDeployment":   kubernetes.GetDeployment{},
		"getHPA":          kubernetes.GetHPA{},
		"getNode":         kubernetes.GetNode{},
		"getPod":          kubernetes.GetPod{},
		"listDeployments": kubernetes.ListDeployments{},
		"listNamespaces":  kubernetes.ListNamespaces{},
		"listPods":        kubernetes.ListPods{},
		"podLogs":         kubernetes.PodLogs{},
		"topPods":         kubernetes.TopPods{},
	}

	k8sWriteExecutors = map[string]IExecutor{
		"deletePod":         kubernetes.DeletePod{},
		"rolloutDeployment": kubernetes.RolloutDeployment{},
	}

	webExecutors = map[string]IExecutor{
		"getUrl": web.GetUrl{},
	}

	commonExecutors = map[string]IExecutor{
		"wait": commonExecutorsPkg.Wait{},
	}
)

func GetExecutors() map[string]IExecutor {
	logger := viper.Get("logger").(logger.ILogger)
	configuration := configuration.LoadConfiguration(logger)

	executors := make(map[string]IExecutor)

	if configuration.Executors.K8S.Read {
		for key, value := range k8sReadExecutors {
			executors[key] = value
		}
	}

	if configuration.Executors.K8S.Write {
		for key, value := range k8sWriteExecutors {
			executors[key] = value
		}
	}

	if configuration.Executors.Web.Enabled {
		for key, value := range webExecutors {
			executors[key] = value
		}
	}

	if configuration.Executors.Common.Enabled {
		for key, value := range commonExecutors {
			executors[key] = value
		}
	}

	return executors
}

type IExecutor interface {
	Exec(executor common.Executor, podName string) string
	GetParams() string
	GetDescription() string
	GetName() string
}
