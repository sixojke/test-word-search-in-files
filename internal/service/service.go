package service

import (
	"github.com/sixojke/internal/config"
	"github.com/sixojke/internal/domain"
	"github.com/sixojke/internal/repository"
)

type Files interface {
	FindFilesWithWord(inp domain.FilesSearchInp) (*domain.FilesSearchOut, error)
}

type Deps struct {
	Repo   *repository.Repository
	Config *config.Config
}

type Service struct {
	Files Files
}

func NewService(deps *Deps) *Service {
	return &Service{Files: NewFilesService(deps.Repo.Cache, deps.Config.Cache)}
}
