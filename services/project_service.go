package services

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/config"
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/Ntrashh/crawlerctl/storage"
	"github.com/Ntrashh/crawlerctl/util"
	"github.com/google/uuid"
	"io/ioutil"
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
	savePath := filepath.Join(fmt.Sprintf("%s/.crawlerctl/projects", config.AppConfig.Path), projectName) // 请根据实际情况修改
	project := &models.Project{
		ProjectName:       projectName,
		VirtualenvName:    virtualEnvName,
		VirtualenvPath:    virtualEnvPath,
		VirtualenvVersion: virtualEnvVersion,
		SavePath:          savePath,
	}
	existingProject, err := p.ProjectStorage.GetByName(projectName)
	if err == nil && existingProject.ID != 0 {
		return fmt.Errorf("项目名称已存在")
	}
	err = p.ProjectStorage.Create(project)
	if err != nil {
		return err
	}
	tempFilePath, err := util.SaveFileToTemp(file)
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(tempFilePath)
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		return err
	}
	err = util.UnzipFile(tempFilePath, savePath)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) GetAllProjects() ([]models.Project, error) {
	allProjects, err := p.ProjectStorage.GetAll()
	if err != nil {
		return nil, err
	}

	// 过滤掉文件夹不存在的项目
	var validProjects = make([]models.Project, 0)
	for _, project := range allProjects {
		// 使用 os.Stat 检查文件夹是否存在
		if _, err := os.Stat(project.SavePath); os.IsNotExist(err) {
			// 文件夹不存在，跳过该项目
			err := p.ProjectStorage.DeleteByID(project.ID)
			if err != nil {
				return nil, err
			}
			continue
		}
		// 文件夹存在，保留该项目
		validProjects = append(validProjects, project)
	}
	return validProjects, nil
}

func (p *ProjectService) DeleteProjectByID(id uint) error {
	project, err := p.ProjectStorage.GetByID(id)
	if err != nil {
		// 如果项目不存在或者获取失败，返回具体的错误
		return fmt.Errorf("无法找到ID为%d的项目: %v", id, err)
	}
	err = p.ProjectStorage.DeleteByID(id)
	if err != nil {
		return err
	}
	path := filepath.Join(fmt.Sprintf("%s/.crawlerctl/projects/%s", config.AppConfig.Path, project.ProjectName))
	err = os.RemoveAll(path) // 删除整个目录或文件
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectService) ProjectsByVersion(virtualenv string) ([]models.Project, error) {
	projects, err := p.ProjectStorage.GetByVersion(virtualenv)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (p *ProjectService) GetFolderTree(folderPath string) ([]map[string]interface{}, error) {
	var tree []map[string]interface{}

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// 检查文件是否没有扩展名，如果没有则跳过该文件
		if !file.IsDir() && filepath.Ext(file.Name()) == "" {
			continue
		}
		node := map[string]interface{}{
			"title":  file.Name(),
			"key":    filepath.Join(folderPath, file.Name()),
			"isLeaf": !file.IsDir(),
		}
		if file.IsDir() {
			// 递归获取子目录
			children, err := p.GetFolderTree(filepath.Join(folderPath, file.Name()))
			if err != nil {
				return nil, err
			}
			node["children"] = children
		}
		tree = append(tree, node)
	}
	return tree, nil
}

func (p *ProjectService) ReadFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
func (p *ProjectService) ProjectById(id uint) (models.Project, error) {
	project, err := p.ProjectStorage.GetByID(id)
	if err != nil {
		return models.Project{}, err
	}
	return *project, nil
}

func (p *ProjectService) SaveFile(filePath, content string) error {
	err := ioutil.WriteFile(filePath, []byte(util.Base64Decode(content)), 0644)
	if err != nil {
		return fmt.Errorf("failed to save code to file: %w", err)
	}
	return nil
}

func (p *ProjectService) ReUploadProject(savePath string, file *multipart.FileHeader) error {

	newUUID, err := uuid.NewUUID()
	var tempDir = fmt.Sprintf("/tmp/%s", newUUID.String())
	defer os.RemoveAll(tempDir)
	err = os.Rename(savePath, tempDir)
	if err != nil {
		return err
	}

	tempFilePath, err := util.SaveFileToTemp(file)
	if err != nil {
		return err
	}
	defer func(name, temp string) {
		_ = os.Remove(name)
		_ = os.Remove(temp)
	}(tempFilePath, tempDir)
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		return err
	}
	err = util.UnzipFile(tempFilePath, savePath)
	if err != nil {
		// 如果解压失败再移动回来
		err = os.Rename(tempDir, savePath)
		return err
	}

	return nil
}
