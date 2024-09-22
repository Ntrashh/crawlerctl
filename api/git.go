package api

import (
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GitHandler struct {
	gitService *services.GitService
}

func NewGitHandler(gitService *services.GitService) *GitHandler {
	return &GitHandler{
		gitService: gitService,
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
	gitConfig, err := g.gitService.CreateGit(gitConfigData.ProjectId, gitConfigData.GitPath, gitConfigData.Username, gitConfigData.Password)
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
