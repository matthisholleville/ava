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

package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/plumbing/transport"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/matthisholleville/ava/pkg/common"
	"github.com/matthisholleville/ava/pkg/knowledge/source/configuration"
	"github.com/matthisholleville/ava/pkg/knowledge/source/local"
	"github.com/matthisholleville/ava/pkg/logger"
)

type Git struct {
	logger        logger.ILogger
	authToken     string
	repositoryURL string
	branch        string
	directory     string
	baseDir       string
}

func (g *Git) Configure(logger logger.ILogger, config configuration.KnowledgeSourceConfiguration) error {
	g.logger = logger
	g.authToken = config.GitAuthToken
	g.repositoryURL = config.GitRepositoryURL
	g.branch = config.GitBranch

	if g.repositoryURL == "" {
		return fmt.Errorf("missing repository")
	}

	baseWorkDir := os.Getenv("WRITE_FS_WORKDIR")
	if baseWorkDir == "" {
		baseWorkDir = "/tmp"
	}

	baseDir := fmt.Sprintf("%s/ava-knowledge-%s", baseWorkDir, common.GenerateRandomString(10))

	g.logger.Debug(fmt.Sprintf("Creating project path %s", g.baseDir))

	err := common.CreateDir(baseDir)
	if err != nil {
		return err
	}

	g.baseDir = baseDir

	g.directory = fmt.Sprintf("%s/%s", baseDir, config.Directory)

	return err
}

func (g *Git) GetFiles() ([]string, error) {
	var files []string

	g.logger.Debug(fmt.Sprintf("Downloading files from %s", g.repositoryURL))

	err := g.cloneRepository()
	if err != nil {
		return files, err
	}

	localSource := &local.Local{}
	localSource.Configure(g.logger, configuration.KnowledgeSourceConfiguration{})

	return localSource.GetFilesFromLocalPath(g.directory)
}

func (g *Git) GetName() string {
	return "git"
}

func (g *Git) CleanUp() error {
	return os.RemoveAll(g.baseDir)
}

func (g *Git) cloneRepository() error {
	var auth transport.AuthMethod
	if g.authToken != "" {
		g.logger.Debug("Using token for authentication")
		auth = &gitHttp.BasicAuth{
			Username: "x-access-token",
			Password: g.authToken,
		}
	}

	loggerWriter := &logger.LoggerWriter{
		Logger: g.logger.Debug,
	}

	g.logger.Debug(fmt.Sprintf("Cloning repository %s", g.repositoryURL))
	cloneOptions := &git.CloneOptions{
		URL:      g.repositoryURL,
		Progress: loggerWriter,
		Auth:     auth,
	}

	if g.branch != "" {
		g.logger.Debug(fmt.Sprintf("Cloning branch %s", g.branch))
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(g.branch)
	}

	g.logger.Debug(fmt.Sprintf("Cloning into %s", g.baseDir))
	_, err := git.PlainClone(g.baseDir, false, cloneOptions)
	return err
}
