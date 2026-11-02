package httpcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/enkidomusic/planne/internal/domain"
	"github.com/enkidomusic/planne/internal/service"
)

type Server struct {
	mux     *http.ServeMux
	service *service.Service
}

func New(service service.Service) *Server {
	s := &Server{
		service: &service,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", s.createTask)
	mux.HandleFunc("GET /tasks", s.getTasks)
	mux.HandleFunc("PUT /tasks/{id}", s.updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", s.deleteTask)

	mux.Handle("GET /", http.FileServer(http.Dir("web")))

	s.mux = mux

	return s
}

func (s *Server) Run(address string) error {
	return http.ListenAndServe(address, s.mux)
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var task *Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("failed to decode req body: %v", err)
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	err = s.service.CreateTask(task.ToDomain())
	if err != nil {
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.service.GetTasks()
	if err != nil {
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	var dtoTasks []*Task

	for _, task := range tasks {
		dtoTasks = append(dtoTasks, DTOFromDomain(task))
	}

	err = json.NewEncoder(w).Encode(&dtoTasks)
	if err != nil {
		log.Printf("failed to encode resp body: %v", err)
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("failed to read task id from path values: %v", err)
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	var task *Task

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("failed to decode req body: %v", err)
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	task.ID = taskID

	err = s.service.UpdateTask(task.ToDomain())
	if err != nil {
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("failed to read task id from path values: %v", err)
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	err = s.service.DeleteTask(taskID)
	if err != nil {
		http.Error(w, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
