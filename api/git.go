package api

import (
	"fmt"
	"github.com/Ntrashh/crawlerctl/config"
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

type GitHandler struct {
	gitService     *services.GitService
	projectService *services.ProjectService
}

func NewGitHandler(gitService *services.GitService, projectService *services.ProjectService) *GitHandler {
	return &GitHandler{
		gitService:     gitService,
		projectService: projectService,
	}
}

func (g GitHandler) GitByProjectIdHandler(c *gin.Context) {
	projectId := c.Param("id")
	projectIdUint64, err := strconv.ParseInt(projectId, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	gitConfig := g.gitService.GetGitByProjectID(int(projectIdUint64))
	if gitConfig == nil {
		SuccessResponse(c, nil)
		return
	}
	SuccessResponse(c, gitConfig)

}

func (g GitHandler) CreateGitConfigHandler(c *gin.Context) {
	var gitConfigData struct {
		ProjectId int    `json:"project_id"`
		GitPath   string `json:"git_path"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}
	err := c.ShouldBind(&gitConfigData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	project, err := g.projectService.ProjectById(uint(gitConfigData.ProjectId))
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "项目不存在")
		return
	}
	projectPath := filepath.Join(fmt.Sprintf("%s/.crawlerctl/projects/%s", config.AppConfig.Path, project.ProjectName))
	gitConfig, err := g.gitService.CreateGit(gitConfigData.ProjectId, gitConfigData.GitPath, gitConfigData.Username, gitConfigData.Password, projectPath)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, gitConfig)
}

func (g GitHandler) RemoteBranchesHandler(c *gin.Context) {
	projectId := c.Param("id")
	projectIdUint64, err := strconv.ParseInt(projectId, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	branches, err := g.gitService.RemoteBranches(int(projectIdUint64))
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, branches)
}

func (g GitHandler) RemoteBranchCommitsHandler(c *gin.Context) {
	projectId := c.Query("id")
	branchName := c.Query("branch_name")

	projectIdUint64, err := strconv.ParseInt(projectId, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	commits, err := g.gitService.GetRemoteBranchCommits(int(projectIdUint64), 5, branchName)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, commits)

}

func (g GitHandler) BranchPullHandler(c *gin.Context) {
	var pullData struct {
		ProjectId  int    `json:"project_id"`
		BranchName string `json:"branch_name"`
	}
	err := c.ShouldBind(&pullData)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	gitConfig, err := g.gitService.BranchPull(pullData.ProjectId, pullData.BranchName)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResponse(c, gitConfig)
}
