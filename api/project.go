package api

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

type ProjectHandler struct {
	ProjectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		ProjectService: projectService,
	}
}

func (p *ProjectHandler) ListProjectsHandler(c *gin.Context) {
	SuccessResponse(c, "")
}

func (p *ProjectHandler) AddProjectHandler(c *gin.Context) {
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

func (p *ProjectHandler) GetAllProjectsHandler(c *gin.Context) {
	projects, err := p.ProjectService.GetAllProjects()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取项目列表失败")
		return
	}
	SuccessResponse(c, projects)
}

func (p *ProjectHandler) DeleteProjectHandler(c *gin.Context) {
	idStr := c.Param("id")
	// 将字符串解析为 uint
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的项目 ID")
		return
	}
	id := uint(idUint64)

	// 调用服务层删除项目
	err = p.ProjectService.DeleteProjectByID(id)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	SuccessResponse(c, "项目已成功删除")
}

func (p *ProjectHandler) ProjectsByVersionHandler(c *gin.Context) {
	version := c.Query("version")
	if version == "" {
		ErrorResponse(c, http.StatusBadRequest, "version不能为空")
		return
	}
	projects, err := p.ProjectService.ProjectsByVersion(version)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	result := true
	if len(projects) == 0 {
		result = false
	}
	SuccessResponse(c, result)
}
