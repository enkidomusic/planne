package repository

import "github.com/enkidomusic/planne/internal/domain"

type Repository interface {
	CreateTask(task *domain.Task) error
	GetTasks() ([]*domain.Task, error)
	UpdateTask(task *domain.Task) error
	DeleteTask(id int) error
}
