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
	store.DeleteAll() //*****************************************brise sve po pokretanju
	for _, User := range users {
		err := store.Insert(User)
		if err != nil {
			log.Fatal(err)
		}
	}
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
	userHandler.Init(r)

	// Dodavanje CORS middleware-a
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Kreiranje HTTP servera sa CORS middleware-om
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: corsHandler(r),
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
