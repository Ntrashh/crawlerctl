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
	projectId := c.Param("projectId")
	projectIdUint64, err := strconv.ParseUint(projectId, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	git := g.gitService.GetGitByProjectID(uint(projectIdUint64))
	if git == nil {
		ErrorResponse(c, http.StatusNotFound, "Not Found")
		return
	}
	SuccessResponse(c, git)

}
