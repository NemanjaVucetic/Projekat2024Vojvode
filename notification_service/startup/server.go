package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notificationService/handler"
	"notificationService/repository"
	"notificationService/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	logger := log.New(os.Stdout, "[notification-api] ", log.LstdFlags)

	store, err := repository.New(logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseSession()
	store.CreateTables()

	err1 := store.DeleteAll()
	if err1 != nil {
		logger.Fatal("Error deleting notifications:", err)
		return
	}

	for _, notification := range notifications {
		err := store.Insert(notification)
		if err != nil {
			logger.Fatal("Error inserting notification:", err)
			return
		}
	}

	notificationService := server.initNotificationService(*store)

	notificationHandler := server.initNotificationHandler(notificationService)

	server.start(notificationHandler)
}

func (server *Server) initNotificationService(store repository.NotificationCassandraStore) *service.NotificationService {
	return service.NewNotificationService(&store)
}

func (server *Server) initNotificationHandler(service *service.NotificationService) *handler.NotificationHandler {
	return handler.NewNotificationHandler(service)
}

func (server *Server) start(orderHandler *handler.NotificationHandler) {
	r := mux.NewRouter()
	orderHandler.Init(r)

	// Set up CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},         // Allow only your frontend domain
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},  // Allowed methods
		AllowedHeaders: []string{"Content-Type", "Authorization"}, // Allowed headers
	})

	// Wrap the router with the CORS middleware
	handler := corsHandler.Handler(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.config.Port),
		Handler: handler, // Use the wrapped router with CORS middleware
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
