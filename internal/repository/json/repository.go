package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/enkidomusic/planne/internal/domain"
)

type JSONRepository struct {
	filepath string
	tasks    map[int]*Task
}

func New(filepath string) (*JSONRepository, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository file: %w", err)
	}
	defer file.Close()

	var tasks []*Task

	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return &JSONRepository{
				filepath: filepath,
				tasks:    make(map[int]*Task),
			}, nil
		}

		return nil, fmt.Errorf("failed to decode json file to tasks struct: %w", err)
	}

	tasksMap := make(map[int]*Task, len(tasks))

	for _, t := range tasks {
		tasksMap[t.ID] = t
	}

	repo := &JSONRepository{
		filepath: filepath,
		tasks:    tasksMap,
	}

	return repo, nil
}

func (r *JSONRepository) Flush() error {
	tasks := make([]*Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		tasks = append(tasks, t)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	dir := filepath.Dir(r.filepath)
	tmp, err := os.CreateTemp(dir, "tasks-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // после успешного Rename файла уже нет, ошибка игнорируется

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmp.Sync(); err != nil { // сбросить на диск до переименования
		tmp.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpName, r.filepath); err != nil {
		return fmt.Errorf("failed to replace json file: %w", err)
	}
	return nil
}

func (r *JSONRepository) CreateTask(task *domain.Task) error {
	// :)
	var maxID int
	for k := range r.tasks {
		if k > maxID {
			maxID = k
		}
	}

	task.ID = maxID + 1
	r.tasks[maxID+1] = DTOFromDomain(task)
	err := r.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}

func (r *JSONRepository) GetTasks() ([]*domain.Task, error) {
	var tasks []*domain.Task

	for _, t := range r.tasks {
		tasks = append(tasks, t.ToDomain())
	}

	return tasks, nil
}

func (r *JSONRepository) UpdateTask(task *domain.Task) error {
	r.tasks[task.ID] = DTOFromDomain(task)
	err := r.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}

func (r *JSONRepository) DeleteTask(id int) error {
	delete(r.tasks, id)
	err := r.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}
