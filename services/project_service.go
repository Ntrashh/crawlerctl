package services

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/Ntrashh/crawlerctl/storage"
	"github.com/Ntrashh/crawlerctl/util"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ProjectService struct {
	ProjectStorage storage.ProjectStorage
}

func NewProjectService(projectStorage storage.ProjectStorage) *ProjectService {
	return &ProjectService{
		ProjectStorage: projectStorage,
	}
}

func (p *ProjectService) AddProjectService(projectName, virtualEnvName, virtualEnvPath, virtualEnvVersion string, file *multipart.FileHeader) error {
	tempFilePath, err := util.SaveFileToTemp(file)
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(tempFilePath)
	destDir := filepath.Join("/opt/crawlerctl/projects", projectName) // 请根据实际情况修改
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = util.UnzipFile(tempFilePath, destDir)
	if err != nil {
		return err
	}
	project := &models.Project{
		ProjectName:       projectName,
		VirtualEnvName:    virtualEnvName,
		VirtualEnvPath:    virtualEnvPath,
		VirtualEnvVersion: virtualEnvVersion,
	}
	existingProject, err := p.ProjectStorage.GetByName(projectName)
	if err == nil && existingProject.ID != 0 {
		return fmt.Errorf("项目名称已存在")
	}
	// 调用仓储层保存项目

	err = p.ProjectStorage.Create(project)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectService) GetAllProjects() ([]models.Project, error) {
	return p.ProjectStorage.GetAll()
}

func (p *ProjectService) DeleteProjectByID(id uint) error {
	// 调用仓储层删除项目
	err := p.ProjectStorage.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}
