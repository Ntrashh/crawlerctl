package storage

import (
	"github.com/Ntrashh/crawlerctl/database"
	"github.com/Ntrashh/crawlerctl/models"
)

type GitStorage interface {
	GetGitByProjectID(projectID int) (*models.Git, error)
	CreateGit(git *models.Git) error
}

type gitStore struct {
}

func (g gitStore) GetGitByProjectID(projectID int) (*models.Git, error) {
	var git models.Git
	result := database.DB.Where("project_id = ?", projectID).First(&git)
	if result.Error != nil {
		return nil, result.Error
	}
	return &git, nil
}

func (g gitStore) CreateGit(git *models.Git) error {
	return database.DB.Create(&git).Error
}

func NewGitStore() GitStorage {
	return &gitStore{}
}
