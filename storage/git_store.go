package storage

import (
	"github.com/Ntrashh/crawlerctl/database"
	"github.com/Ntrashh/crawlerctl/models"
)

type GitStorage interface {
	GetGitByProjectID(projectID uint) (*models.Git, error)
}

type gitStore struct {
}

func (g gitStore) GetGitByProjectID(projectID uint) (*models.Git, error) {
	var git models.Git
	result := database.DB.Where("projectID = ?", projectID).First(&git)
	if result.Error != nil {
		return nil, result.Error
	}
	return &git, nil
}

func NewGitStore() GitStorage {
	return &gitStore{}
}
