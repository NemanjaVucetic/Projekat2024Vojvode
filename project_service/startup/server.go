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
	store.DeleteAll()
	for _, Project := range projects {
		err := store.Insert(Project)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initProjectService(store repository.ProjectMongoDBStore, userClient client.Client) *service.ProjectService {
	return service.NewProjectService(store, userClient)
}

func (server *Server) initProjectHandler(service *service.ProjectService) *handler.ProjectsHandler {
	return handler.NewProjectsHandler(service)
}

func (server *Server) start(orderHandler *handler.ProjectsHandler) {
	r := mux.NewRouter()
	orderHandler.Init(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: r,
	}

	wait := time.Second * 15
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server %s", err)
	}
	log.Println("server gracefully stopped")
}
