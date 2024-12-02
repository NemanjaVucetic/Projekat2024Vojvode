package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"taskService/handler"
	"taskService/repository"
	"taskService/service"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[task-store] ", log.LstdFlags)

	taskStore, err := repository.NewTaskRepo(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer taskStore.Disconnect(timeoutContext)

	tasksService := service.NewTaskService(taskStore)

	//Initialize the handler and inject said logger
	taskHandler := handler.NewTaskHandler(logger, tasksService)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()

	pingRouter := router.Methods(http.MethodGet).Subrouter()
	pingRouter.HandleFunc("/ping", healthCheckHandler)

	patchRouterTask := router.Methods(http.MethodPatch).Subrouter()
	patchRouterTask.HandleFunc("/tasks", taskHandler.CreateTask)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")

}
