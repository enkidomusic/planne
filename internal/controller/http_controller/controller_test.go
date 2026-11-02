package httpcontroller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/enkidomusic/planne/internal/domain"
	"github.com/enkidomusic/planne/internal/service"
)

type recordingRepository struct {
	created []*domain.Task
	updated []*domain.Task
}

func (r *recordingRepository) CreateTask(task *domain.Task) error {
	r.created = append(r.created, task)
	return nil
}

func (r *recordingRepository) UpdateTask(task *domain.Task) error {
	r.updated = append(r.updated, task)
	return nil
}

func (r *recordingRepository) GetTasks() ([]*domain.Task, error) {
	panic("unexpected GetTasks call")
}

func (r *recordingRepository) DeleteTask(int) error {
	panic("unexpected DeleteTask call")
}

func TestTaskHandlersRejectNull(t *testing.T) {
	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/tasks"},
		{http.MethodPut, "/tasks/42"},
	} {
		t.Run(tc.method, func(t *testing.T) {
			for _, body := range []string{"null", " \n null \t\n"} {
				t.Run(body, func(t *testing.T) {
					repo := &recordingRepository{}
					server := New(*service.New(repo))
					req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(body))
					req.Header.Set("Content-Type", "application/json")
					response := httptest.NewRecorder()

					// Serve through the mux so PUT receives its path ID. Any panic fails the test.
					server.mux.ServeHTTP(response, req)

					if response.Code != http.StatusBadRequest {
						t.Fatalf("status = %d, want %d; body = %q", response.Code, http.StatusBadRequest, response.Body.String())
					}
					if len(repo.created) != 0 || len(repo.updated) != 0 {
						t.Fatalf("null request reached repository: created=%d, updated=%d", len(repo.created), len(repo.updated))
					}
				})
			}
		})
	}
}

func TestTaskHandlersAcceptTaskObject(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		t.Run(method, func(t *testing.T) {
			repo := &recordingRepository{}
			server := New(*service.New(repo))
			path := "/tasks"
			if method == http.MethodPut {
				path = "/tasks/42"
			}
			req := httptest.NewRequest(method, path, strings.NewReader(`{"id":99,"description":"Recurring task","repeat_count":3,"repeat_interval":2}`))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			server.mux.ServeHTTP(response, req)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body = %q", response.Code, http.StatusOK, response.Body.String())
			}
			var task *domain.Task
			if method == http.MethodPost {
				if len(repo.created) != 1 || len(repo.updated) != 0 {
					t.Fatalf("repository calls: created=%d, updated=%d", len(repo.created), len(repo.updated))
				}
				task = repo.created[0]
			} else {
				if len(repo.updated) != 1 || len(repo.created) != 0 {
					t.Fatalf("repository calls: created=%d, updated=%d", len(repo.created), len(repo.updated))
				}
				task = repo.updated[0]
				if task.ID != 42 {
					t.Errorf("task ID = %d, want path ID 42", task.ID)
				}
			}
			if task.Description != "Recurring task" || task.RepeatCount == nil || *task.RepeatCount != 3 || task.RepeatInterval == nil || *task.RepeatInterval != 2 {
				t.Errorf("unexpected task received by repository: %+v", task)
			}
		})
	}
}
