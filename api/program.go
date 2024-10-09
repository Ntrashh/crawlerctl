package api

import (
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProgramHandler struct {
	programService *services.ProgramService
}

func NewProgramHandler(programService *services.ProgramService) *ProgramHandler {
	return &ProgramHandler{
		programService: programService,
	}
}

func (p *ProgramHandler) CreateProgram(c *gin.Context) {
	var req struct {
		ProgramName  string `json:"program_name" binding:"required"`
		ProjectId    int    `json:"project_id" binding:"required"`
		StartCommand string `json:"start_command" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := p.programService.AddProgramService(req.ProgramName, req.StartCommand, uint(req.ProjectId))
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, nil)
}

func (p *ProgramHandler) GetPrograms(c *gin.Context) {
	programs, err := p.programService.GetAllPrograms()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(c, programs)
}
