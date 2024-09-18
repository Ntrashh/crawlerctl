package storage

import (
	"github.com/Ntrashh/crawlerctl/database"
	"github.com/Ntrashh/crawlerctl/models"
)

type ProjectStorage interface {
	Create(project *models.Project) error
	GetByName(name string) (*models.Project, error)
	// 根据需要添加更多方法
}
type projectStorage struct {
	// 可以添加依赖，例如数据库连接
}

func (p projectStorage) Create(project *models.Project) error {
	return database.DB.Create(project).Error
}

func (p projectStorage) GetByName(name string) (*models.Project, error) {
	var project models.Project
	result := database.DB.Where("project_name = ?", name).FirstOrInit(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func NewProjectStorage() ProjectStorage {
	return &projectStorage{}
}
