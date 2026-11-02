package app

import (
	httpcontroller "github.com/enkidomusic/planne/internal/controller/http_controller"
	jsonRepo "github.com/enkidomusic/planne/internal/repository/json"
	"github.com/enkidomusic/planne/internal/service"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	repo, err := jsonRepo.New("./data/tasks.json")
	if err != nil {
		panic(err)
	}

	svc := service.New(repo)

	srv := httpcontroller.New(*svc)

	srv.Run("localhost:8080")
}
