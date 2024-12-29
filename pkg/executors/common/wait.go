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

package common

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/matthisholleville/ava/pkg/common"
)

type Wait struct {
	Time string `json:"time"`
}

func (Wait) Exec(e common.Executor, jsonString string) string {
	// execute a check on the url and return the status code
	var waitInfo Wait
	err := json.Unmarshal([]byte(jsonString), &waitInfo)
	if err != nil {
		return "Error while retrieving the wait parameter:" + err.Error()
	}
	
	// Convert the wait time to an integer
	waitSeconds, err := strconv.Atoi(waitInfo.Time)
	if err != nil {
		return "Error while parsing the wait time: " + err.Error()
	}

	// Wait for the specified time
	time.Sleep(time.Duration(waitSeconds) * time.Second)
	return "Waited for " + waitInfo.Time + " seconds"
}
