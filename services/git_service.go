package services

import (
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/Ntrashh/crawlerctl/storage"
)

type GitService struct {
	GitStorage storage.GitStorage
}

// NewEnvService 创建一个新的 EnvService 实例
func NewGitService(gitStorage storage.GitStorage) *GitService {
	return &GitService{
		GitStorage: gitStorage,
	}
}

func (g GitService) GetGitByProjectID(projectID uint) *models.Git {
	git, err := g.GitStorage.GetGitByProjectID(projectID)
	if err != nil {
		return nil
	}
	return git
}
