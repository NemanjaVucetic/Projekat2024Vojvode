package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userService/handler"
	"userService/repository"
	"userService/service"

	"github.com/gorilla/handlers" // Import gorilla handlers za CORS
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

	userStore := server.initUserStore(mongoClient)

	userService := server.initUserService(*userStore)

	userHandler := server.initUserHandler(*userService)

	server.start(userHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := repository.GetClient(server.config.UserDBHost, server.config.UserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initUserStore(client *mongo.Client) *repository.UserMongoDBStore {
	store := repository.NewUserMongoDBStore(client)
	//store.DeleteAll() //*****************************************brise sve po pokretanju
	//for _, User := range users {
	//err := store.Insert(User)
	//if err != nil {
	//log.Fatal(err)
	//}
	//}
	return store
}

func (server *Server) initUserService(store repository.UserMongoDBStore) *service.UserService {
	return service.NewUserService(store)
}

func (server *Server) initUserHandler(service service.UserService) *handler.UsersHandler {
	return handler.NewUsersHandler(service)
}

func (server *Server) start(userHandler *handler.UsersHandler) {
	r := mux.NewRouter()

	// Add a handler for preflight requests
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent) // Return HTTP 204
	})

	// Initialize the userHandler
	userHandler.Init(r)

	// Add CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	// Wrap the router with CORS middleware
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: corsHandler(r),
	}

	// Graceful shutdown logic with a timeout of 15 seconds
	wait := time.Second * 15
	go func() {
		log.Printf("Starting server on port %s", server.config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Handle graceful shutdown
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
