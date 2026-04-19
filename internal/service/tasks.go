package service

import (
	"lefree/go-todo/internal/models"
	"log/slog"
)

type Repository interface {
	Save(dto *models.Task) (string, error)
	GetAll() ([]*models.Task, error)
	GetById(id string) (*models.Task, error)
	Patch(dto *models.Task) (*models.Task, error)
	Delete(id string) error
}

type TaskService struct {
	logger *slog.Logger
	repo   Repository
}

func NewTaskService(log *slog.Logger, repository Repository) *TaskService {
	return &TaskService{
		logger: log,
		repo:   repository,
	}
}

func (s *TaskService) AddTask(dto *models.Task) (string, error) {
	return s.repo.Save(dto)
}

func (s *TaskService) GetAll() ([]*models.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) GetById(id string) (*models.Task, error) {
	return s.repo.GetById(id)
}

func (s *TaskService) Patch(dto *models.Task) (*models.Task, error) {
	return s.repo.Patch(dto)
}

func (s *TaskService) Delete(id string) error {
	return s.repo.Delete(id)
}
