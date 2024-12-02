package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"taskService/domain"
	"taskService/service"
)

type KeyProduct struct{}

type TaskHandler struct {
	logger  *log.Logger
	service *service.TaskService
}

// Injecting the logger makes this code much more testable.
func NewTaskHandler(l *log.Logger, s *service.TaskService) *TaskHandler {
	return &TaskHandler{logger: l, service: s}
}

func (p *TaskHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		p.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

func (p *TaskHandler) MiddlewareProjectDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		task := &domain.Task{}
		err := task.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			p.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, task)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (p *TaskHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Optionally, perform any necessary checks here (database connection, etc.)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}

func (p *TaskHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	// Extract the task from the request body
	task := &domain.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
		p.logger.Println(err)
		return
	}

	// Creating the task using the service
	ctx := r.Context() // Get the context from the HTTP request
	createdTask, err := p.service.Create(ctx, task)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responding with the created task as JSON
	rw.WriteHeader(http.StatusCreated)
	err = createdTask.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to serialize task", http.StatusInternalServerError)
	}
}
