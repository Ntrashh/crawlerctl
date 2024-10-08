package storage

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/database"
	"github.com/Ntrashh/crawlerctl/models"
)

type ProjectStorage interface {
	Create(project *models.Project) error
	GetByName(name string) (*models.Project, error)
	GetByVersion(version string) ([]models.Project, error)
	GetByID(id uint) (*models.Project, error)
	GetAll() ([]models.Project, error)
	DeleteByID(id uint) error
}
type projectStorage struct {
	// 可以添加依赖，例如数据库连接
}

func (p projectStorage) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	result := database.DB.Where("ID = ?", id).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (p projectStorage) GetByVersion(version string) ([]models.Project, error) {
	var projects []models.Project
	result := database.DB.Where("virtualenv_version = ?", version).Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (p projectStorage) DeleteByID(id uint) error {
	result := database.DB.Unscoped().Delete(&models.Project{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("项目未找到")
	}
	return nil
}

func (p projectStorage) GetAll() ([]models.Project, error) {
	var projects = make([]models.Project, 0)
	result := database.DB.Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
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
