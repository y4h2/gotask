package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

// App defines app server
type App struct {
	server   *http.Server
	listener net.Listener
}

// Start starts server
func (app *App) Start() error {
	listener, err := net.Listen("tcp", app.server.Addr)
	if err != nil {
		return err
	}
	app.listener = listener

	go app.server.Serve(app.listener)

	return nil
}

// Shutdown shuts down the server
func (app *App) Shutdown() {
	app.server.Shutdown(context.Background())
}

// GetPort returns http server's binding port
func (app *App) GetPort() int {
	return app.listener.Addr().(*net.TCPAddr).Port
}

// NewApp is the constructor of App
func NewApp(port int) *App {
	router := mux.NewRouter()

	app := &App{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
	}
	return app
}

func Registery(router *mux.Router) {

	// taskServer := NewTaskService()

	httpServer := NewHttpServer()

	router.HandleFunc("/ping", httpServer.pingHandler)
	router.HandleFunc("/task", httpServer.ListTasksHandler).
		Methods("GET")
	router.HandleFunc("/task/{task-id}", httpServer.CheckTaskHandler).
		Methods("GET")
	router.HandleFunc("/task/{task-id}", httpServer.StartTaskHandler).
		Methods("DELETE")
	router.HandleFunc("/task/{task-id}", httpServer.CancelTaskHandler).
		Methods("POST")
}

type HttpServer struct {
	taskService *TaskService
}

func (h HttpServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h HttpServer) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h HttpServer) CheckTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["task-id"]

	w.WriteHeader(http.StatusOK)
}

func (h HttpServer) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["task-id"]

	w.WriteHeader(http.StatusOK)

}

func (h HttpServer) CancelTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["task-id"]

	w.WriteHeader(http.StatusOK)
}

func NewHttpServer() *HttpServer {
	return &HttpServer{}
}
