package services

import (
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/Ntrashh/crawlerctl/storage"
)

type ProgramService struct {
	programStore storage.ProgramStore
}

func NewProgramService(projectStorage storage.ProgramStore) *ProgramService {
	return &ProgramService{
		programStore: projectStorage,
	}
}

func (p *ProgramService) AddProgramService(programName, startCommand string, projectId uint) error {
	programMode := models.Program{
		Name:         programName,
		StartCommand: startCommand,
		ProjectID:    projectId,
	}
	err := p.programStore.Create(programMode)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProgramService) GetAllPrograms() ([]models.Program, error) {
	programs, err := p.programStore.GetAll()
	if err != nil {
		return nil, err
	}
	return programs, nil
}
