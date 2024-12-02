package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projectService/handler"
	"projectService/repository"
	"projectService/service"
	"syscall"
	"time"

	"projectService/client"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	config *Config
}

func NewServer(config1 *Config) *Server {
	return &Server{
		config: config1,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Printf("error closing db: %s\n", err)
		}
	}(mongoClient, context.Background())

	projectStore := server.initProjectStore(mongoClient)

	userClient := server.initUserClient()

	projectService := server.initProjectService(*projectStore, userClient)

	projectHandler := server.initProjectHandler(projectService)

	server.start(projectHandler)
}

func (server *Server) initUserClient() client.Client {
	return client.NewClient(server.config.UserHost, server.config.UserPort)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := repository.GetClient(server.config.ProjectDBHost, server.config.ProjectDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initProjectStore(client *mongo.Client) *repository.ProjectMongoDBStore {
	store := repository.NewProjectMongoDBStore(client)
	//store.DeleteAll()
	//for _, Project := range projects {
	//err := store.Insert(Project)
	//if err != nil {
	//log.Fatal(err)
	//}
	//}
	return store
}

func (server *Server) initProjectService(store repository.ProjectMongoDBStore, userClient client.Client) *service.ProjectService {
	return service.NewProjectService(store, userClient)
}

func (server *Server) initProjectHandler(service *service.ProjectService) *handler.ProjectsHandler {
	return handler.NewProjectsHandler(service)
}

func (server *Server) start(projectHandler *handler.ProjectsHandler) {
	r := mux.NewRouter()

	// Add a handler for preflight requests
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent) // Return HTTP 204
	})

	// Initialize the handler
	projectHandler.Init(r)

	// Add CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	// Wrap the router
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: corsHandler(r),
	}

	// Start the server
	wait := time.Second * 15
	go func() {
		log.Printf("Starting server on port %s", server.config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Graceful shutdown logic
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %s", err)
	}
	log.Println("Server gracefully stopped")
}
