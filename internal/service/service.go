package service

import (
	"log"

	"github.com/enkidomusic/planne/internal/domain"
	"github.com/enkidomusic/planne/internal/repository"
)

type Service struct {
	repository repository.Repository
}

func New(repository repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateTask(task *domain.Task) error {
	if task.Description == "" {
		return domain.ErrInvalidInput
	}

	if (task.RepeatCount == nil && task.RepeatInterval != nil) || (task.RepeatCount != nil && task.RepeatInterval == nil) {
		return domain.ErrInvalidInput
	}

	err := s.repository.CreateTask(task)
	if err != nil {
		log.Printf("failed to create task in repo: %v", err)
		return domain.ErrInternal
	}

	return nil
}

func (s *Service) GetTasks() ([]*domain.Task, error) {
	tasks, err := s.repository.GetTasks()
	if err != nil {
		log.Printf("failed to get tasks: %v", err)
		return nil, domain.ErrInternal
	}

	return tasks, nil
}

func (s *Service) UpdateTask(task *domain.Task) error {
	if task.Description == "" {
		return domain.ErrInvalidInput
	}

	if (task.RepeatCount == nil && task.RepeatInterval != nil) || (task.RepeatCount != nil && task.RepeatInterval == nil) {
		return domain.ErrInvalidInput
	}

	err := s.repository.UpdateTask(task)
	if err != nil {
		log.Printf("failed to update task: %v", err)
		return domain.ErrInternal
	}

	return nil
}

func (s *Service) DeleteTask(id int) error {
	err := s.repository.DeleteTask(id)
	if err != nil {
		log.Printf("failed to delete task: %v", err)
		return domain.ErrInternal
	}

	return nil
}
