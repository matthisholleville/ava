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
	"github.com/matthisholleville/ava/pkg/common"
	commonExecutorsPkg "github.com/matthisholleville/ava/pkg/executors/common"
	"github.com/matthisholleville/ava/pkg/executors/kubernetes"
	"github.com/matthisholleville/ava/pkg/executors/web"
)

var (
	k8sExecutors = map[string]IExecutor{
		"getPod":    kubernetes.GetPod{},
		"podLogs":   kubernetes.LogsPod{},
		"deletePod": kubernetes.DeletePod{},
	}

	webExecutors = map[string]IExecutor{
		"getUrl": web.GetUrl{},
	}

	commonExecutors = map[string]IExecutor{
		"wait": commonExecutorsPkg.Wait{},
	}
)

func GetExecutors() map[string]IExecutor {
	executors := make(map[string]IExecutor)
	for key, value := range k8sExecutors {
		executors[key] = value
	}

	for key, value := range webExecutors {
		executors[key] = value
	}

	for key, value := range commonExecutors {
		executors[key] = value
	}

	return executors
}

type IExecutor interface {
	Exec(executor common.Executor, podName string) string
}
