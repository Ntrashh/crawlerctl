package api

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type ProjectHandler struct {
	ProjectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		ProjectService: projectService,
	}
}

func (p *ProjectHandler) ListProjects(c *gin.Context) {
	SuccessResponse(c, "")
}

func (p *ProjectHandler) AddProject(c *gin.Context) {
	// 获取表单参数
	projectName := c.PostForm("project_name")
	virtualEnvName := c.PostForm("virtualenv_name")
	virtualEnvPath := c.PostForm("virtualenv_path")
	virtualEnvVersion := c.PostForm("virtualenv_version")
	if projectName == "" || virtualEnvName == "" || virtualEnvPath == "" || virtualEnvVersion == "" {
		ErrorResponse(c, http.StatusBadRequest, "缺少必要的参数")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无法获取上传的文件")
		return
	}

	// 验证文件类型是否为 ZIP
	if filepath.Ext(file.Filename) != ".zip" {
		ErrorResponse(c, http.StatusBadRequest, "只能上传 ZIP 文件")
		return
	}
	err = p.ProjectService.AddProjectService(projectName, virtualEnvName, virtualEnvPath, virtualEnvVersion, file)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("解压文件失败：%v", err))
		return
	}
	SuccessResponse(c, "文件上传并解压成功")
}
