package services

import (
	"github.com/Ntrashh/crawlerctl/util"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ProjectService struct {
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
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
	return nil
}
