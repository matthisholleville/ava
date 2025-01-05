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

package types

const VECTOR_STORE_NAME = "ava-sre-agent"
const ASSISTANT_NAME = "ava-sre-agent"
const ASSISTANT_INSTRUCTIONS = `
You are an expert Site Reliability Engineer, tasked with helping
the SRE team respond to and resolve incidents. Please use vector
search to find the most relevant information to help you respond
and take action.

If you are presented with a question that does not seem like it
could be related to infrastructure, begin your response with a polite
reminder that your primary responsibilities are to help with incident
response, before fully answering the question to the best of your ability.
`

var (
	// ANALYSE_AND_FIX_PROMPT = `
	// Using the provided runbooks, follow the outlined steps to execute the necessary functions and resolve the issue.
	// You must respond to the user in %s with a detailed explanation of the steps that allowed you to understand and fix the problem.

	// The problem: %s
	// `

	ANALYSE_AND_FIX_PROMPT = `
	You need to analyze a problem and attempt to resolve it. Break down your reasoning as follows:

	1. Check if a runbook exists in the file search. If it does, follow the instructions in that runbook to address the issue and try to solve it using the functions available to you.
	2. If no runbook is available, inform the user that you will do your best to assist them using your general knowledge and the functions at your disposal.

	You must respond to the user in %s with a detailed explanation of the steps that allowed you to understand and fix the problem.

	The problem: %s
	`

	ANALYSE_PROMPT = `
	Using the provided runbooks, help the user understand the problem and provide a detailed explanation of the steps he can follow to fix the problem.

	You must respond to the user in %s.
	
	The problem: %s`
)
