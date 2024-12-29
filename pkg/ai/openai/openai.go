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

package openai

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/matthisholleville/ava/pkg/ai/types"
	"github.com/matthisholleville/ava/pkg/common"
	"github.com/matthisholleville/ava/pkg/executors"
	"github.com/matthisholleville/ava/pkg/logger"
	"github.com/sashabaranov/go-openai"
)

const openaiClientName = "openai"

type OpenAIClient struct {
	client        *openai.Client
	Configuration Configuration
	ctx           context.Context
	logger        logger.ILogger
}

type Configuration struct {
	VectorID    string `json:"vectorID"`
	AssistantID string `json:"assistantID"`
	APIKey      string `json:"apiKey"`
}

func (c *OpenAIClient) Configure(logger logger.ILogger) error {

	c.client = openai.NewClient(c.Configuration.APIKey)
	c.ctx = context.Background()
	c.logger = logger

	return nil
}

func (c *OpenAIClient) getVectorStore() error {
	c.logger.Debug("Checking if the vector store exists")
	vectors, err := c.client.ListVectorStores(c.ctx, openai.Pagination{})
	if err != nil {
		return err
	}

	for _, vector := range vectors.VectorStores {
		// We are looking for the vector store named "ava-sre-agent"
		if vector.Name == types.VECTOR_STORE_NAME {
			c.Configuration.VectorID = vector.ID
			return nil
		}
	}

	return fmt.Errorf("vector store not found")
}

func (c *OpenAIClient) Analyze(text, language string, threadID string, executorConfig common.Executor) (string, error) {
	var response string

	c.logger.Info(fmt.Sprintf("Analyzing the message: %s", text))

	inputMessage := fmt.Sprintf(types.ANALYSE_PROMPT, language, text)

	c.logger.Debug("Creating a message")
	_, err := c.createMessage(threadID, inputMessage)
	if err != nil {
		return response, err
	}

	c.logger.Debug(fmt.Sprintf("Debug link: https://platform.openai.com/playground/assistants?assistant=%s&thread=%s", c.Configuration.AssistantID, threadID))

	c.logger.Debug("Creating a run")
	run, err := c.createRun(threadID)
	if err != nil {
		return response, err
	}

	c.logger.Debug("Watching the run")
	_, err = c.watchRun(executorConfig, threadID, run.ID)
	if err != nil {
		return response, err
	}

	messages, err := c.listThreadMessage(threadID, run.ID)
	if err != nil {
		return response, err
	}

	response = messages.Messages[0].Content[0].Text.Value

	return response, nil
}

func (c *OpenAIClient) listThreadMessage(threadId, runID string) (messages openai.MessagesList, err error) {
	messages, err = c.client.ListMessage(c.ctx, threadId, nil, nil, nil, nil, &runID)
	return messages, err
}

func (c *OpenAIClient) retrieveRun(threadId, runId string) (*openai.Run, error) {
	response, err := c.client.RetrieveRun(c.ctx, threadId, runId)
	return &response, err
}

func (c *OpenAIClient) watchRun(e common.Executor, threadId, runId string) (*openai.Run, error) {
	completed := false

	c.logger.Info("Analysis in progress...")

	for !completed {

		time.Sleep(1 * time.Second)

		c.logger.Debug("Retrieving the run")
		run, err := c.retrieveRun(threadId, runId)
		if err != nil {
			return nil, err
		}

		if run.Status == "failed" {
			return nil, errors.New("run failed")
		}

		c.logger.Debug(fmt.Sprintf("Run status: %s", run.Status))
		if run.Status == "requires_action" {
			functionsToBeCalled := run.RequiredAction.SubmitToolOutputs

			outputs := []openai.ToolOutput{}

			executors := executors.GetExecutors()

			for _, f := range functionsToBeCalled.ToolCalls {

				c.logger.Info(fmt.Sprintf("Execution of the function : %s", f.Function.Name))
				outputs = append(outputs, openai.ToolOutput{
					ToolCallID: f.ID,
					Output:     executors[f.Function.Name].Exec(e, f.Function.Arguments),
				})
			}

			c.logger.Debug("Submitting the tool outputs")
			run, err = c.submitToolOutputs(threadId, run.ID, outputs)
			if err != nil {
				return nil, err
			}
			return c.watchRun(e, threadId, run.ID)
		}

		if run.Status == "completed" {
			return run, nil
		}

	}

	return nil, errors.ErrUnsupported
}

func (c *OpenAIClient) submitToolOutputs(threadId, runId string, outputs []openai.ToolOutput) (*openai.Run, error) {
	run, err := c.client.SubmitToolOutputs(c.ctx, threadId, runId, openai.SubmitToolOutputsRequest{
		ToolOutputs: outputs,
	})
	return &run, err
}

func (c *OpenAIClient) createMessage(threadId, content string) (*openai.Message, error) {
	message, err := c.client.CreateMessage(c.ctx, threadId, openai.MessageRequest{
		Role:    "user",
		Content: content,
	})
	return &message, err
}

func (c *OpenAIClient) createRun(threadId string) (*openai.Run, error) {
	run, err := c.client.CreateRun(c.ctx, threadId, openai.RunRequest{
		AssistantID: c.Configuration.AssistantID,
	})
	return &run, err
}

func (c *OpenAIClient) CreateThread() (*string, error) {
	thread, err := c.client.CreateThread(c.ctx, openai.ThreadRequest{})
	return &thread.ID, err
}

func (c *OpenAIClient) ConfigureKnowledge(logger logger.ILogger) error {
	err := c.Configure(logger)
	if err != nil {
		return err
	}

	err = c.getVectorStore()
	if err != nil {
		c.logger.Warn("Vector store not found")
	}

	if c.Configuration.VectorID == "" {
		// If the vector store does not exist, we create it
		c.logger.Debug("Creating the vector store")
		vector, err := c.client.CreateVectorStore(c.ctx, openai.VectorStoreRequest{
			Name: types.VECTOR_STORE_NAME,
		})
		if err != nil {
			return err
		}

		c.Configuration.VectorID = vector.ID
	}

	return nil
}

func (c *OpenAIClient) ConfigureAssistant(logger logger.ILogger) error {
	err := c.Configure(logger)
	if err != nil {
		return err
	}
	c.logger.Debug("Checking if the assistant exists")
	assistants, err := c.client.ListAssistants(c.ctx, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	assistantName := types.ASSISTANT_NAME
	assistantInstructions := types.ASSISTANT_INSTRUCTIONS

	for _, assistant := range assistants.Assistants {
		if assistant.Name != nil && *assistant.Name == assistantName {
			c.logger.Debug("Assistant found")
			c.Configuration.AssistantID = assistant.ID
		}
	}

	err = c.getVectorStore()
	if err != nil {
		c.logger.Warn("Vector store not found")
	}

	if c.Configuration.AssistantID == "" {
		// If the assistant does not exist, we create it
		c.logger.Debug("Creating the assistant")
		assistant, err := c.client.CreateAssistant(c.ctx, openai.AssistantRequest{
			Name:  &assistantName,
			Model: openai.GPT4o,
		})
		if err != nil {
			return err
		}

		c.Configuration.AssistantID = assistant.ID
	}

	tools := []openai.AssistantTool{
		{
			Type: openai.AssistantToolTypeFileSearch,
		},
	}
	tools = append(tools, Functions...)
	c.logger.Debug("Modifying the assistant")
	_, err = c.client.ModifyAssistant(c.ctx, c.Configuration.AssistantID, openai.AssistantRequest{
		Instructions: &assistantInstructions,
		Tools:        tools,
		ToolResources: &openai.AssistantToolResource{
			FileSearch: &openai.AssistantToolFileSearch{
				VectorStoreIDs: []string{c.Configuration.VectorID},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *OpenAIClient) GetName() string {
	return openaiClientName
}

func (c *OpenAIClient) UploadFiles(files []string) error {

	c.logger.Debug("Listing already uploaded files")
	alreadyUploadedFiles, err := c.client.ListFiles(c.ctx)
	if err != nil {
		return err
	}

	uploadedFiles := []string{}
	c.logger.Debug("Uploading files")
	for _, file := range files {

		// Check if the file is already uploaded
		c.logger.Debug(fmt.Sprintf("Checking if the file %s is already uploaded", file))
		fileAlreadyExist, fileID := c.isFileAlreadyUploaded(file, alreadyUploadedFiles.Files)

		// If the file is already uploaded, we delete it
		if fileAlreadyExist {
			c.logger.Debug(fmt.Sprintf("The file %s is already uploaded, deleting it", file))
			err = c.client.DeleteFile(c.ctx, fileID)
			if err != nil {
				return err
			}
		}

		// We upload the file
		c.logger.Debug(fmt.Sprintf("Uploading the file %s", file))
		uploadedFile, err := c.client.CreateFile(c.ctx, openai.FileRequest{
			FileName: filepath.Base(file),
			FilePath: file,
			Purpose:  "assistants",
		})
		if err != nil {
			return err
		}

		// We keep track of the uploaded files
		uploadedFiles = append(uploadedFiles, uploadedFile.ID)
	}

	// We upload the vector file
	c.logger.Debug("Uploading the vector files")
	err = c.UploadVectorFiles(uploadedFiles)

	return err
}

func (c *OpenAIClient) Purge() error {
	fileIDs, err := c.listVectorFiles()
	if err != nil {
		return err
	}
	for _, fileID := range fileIDs {
		c.logger.Debug(fmt.Sprintf("Deleting file %s", fileID))
		err := c.client.DeleteFile(c.ctx, fileID)
		if err != nil {
			c.logger.Warn(fmt.Sprintf("Error deleting file %s: %s", fileID, err.Error()))
			continue
		}
	}

	// Delete the vector store
	c.logger.Debug("Deleting the vector store")
	_, err = c.client.DeleteVectorStore(c.ctx, c.Configuration.VectorID)

	return err
}

func (c *OpenAIClient) listVectorFiles() ([]string, error) {
	c.logger.Debug("Listing vector files")
	files, err := c.client.ListVectorStoreFiles(c.ctx, c.Configuration.VectorID, openai.Pagination{})
	if err != nil {
		return nil, err
	}
	fileIDs := make([]string, len(files.VectorStoreFiles))
	c.logger.Debug(fmt.Sprintf("Found %d vector files", len(files.VectorStoreFiles)))
	for _, file := range files.VectorStoreFiles {
		fileIDs = append(fileIDs, file.ID)
	}

	return fileIDs, nil
}

func (c *OpenAIClient) isFileAlreadyUploaded(file string, uploadedFiles []openai.File) (bool, string) {
	for _, uploadedFile := range uploadedFiles {
		if file == uploadedFile.FileName {
			return true, uploadedFile.ID
		}
	}
	return false, ""
}

func (c *OpenAIClient) UploadVectorFiles(fileIDs []string) error {
	_, err := c.client.CreateVectorStoreFileBatch(
		c.ctx,
		c.Configuration.VectorID,
		openai.VectorStoreFileBatchRequest{
			FileIDs: fileIDs,
		},
	)
	return err

}
