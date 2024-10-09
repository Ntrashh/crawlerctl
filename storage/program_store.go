package storage

import (
	"github.com/Ntrashh/crawlerctl/database"
	"github.com/Ntrashh/crawlerctl/models"
)

type ProgramStore interface {
	Create(program models.Program) error
	GetAll() ([]models.Program, error)
}

type programStore struct {
}

func NewProgramStore() ProgramStore {
	return &programStore{}
}

func (s programStore) Create(program models.Program) error {
	return database.DB.Create(program).Error
}

func (s programStore) GetAll() ([]models.Program, error) {
	var programs = make([]models.Program, 0)
	result := database.DB.Preload("Project").Find(&programs)
	if result.Error != nil {
		return nil, result.Error
	}
	return programs, nil

}
