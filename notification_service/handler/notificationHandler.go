package handler

import (
	"encoding/json"
	"net/http"
	"notificationService/service"
	"time"

	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

func (handler *NotificationHandler) Init(r *mux.Router) {
	r.HandleFunc("/", handler.GetAllNotifications).Methods("GET")
	r.HandleFunc("/{id}", handler.GetNotificationByID).Methods("GET")
	r.HandleFunc("/discard/{user_id}/{id}", handler.SetNotificationStatusDiscarded).Methods("PUT")
	r.HandleFunc("/byUserId/{userId}", handler.GetNotificationsByUserID).Methods("GET")
	http.Handle("/", r)
}

func (handler *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := handler.service.GetAll()
	if err != nil {
		http.Error(w, "Unable to fetch notifications", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, "Unable to encode notifications to JSON", http.StatusInternalServerError)
	}
}

func (handler *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	notification, err := handler.service.Get(id)
	if err != nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(notification); err != nil {
		http.Error(w, "Unable to encode Notification to JSON", http.StatusInternalServerError)
	}
}

func (handler *NotificationHandler) GetNotificationsByUserID(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the route parameters
	vars := mux.Vars(r)
	userId, ok := vars["userId"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Call the service method to get notifications for the user
	notifications, err := handler.service.GetByUserId(userId)
	if err != nil {
		http.Error(w, "Failed to fetch notifications for user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the notifications as JSON and send the response
	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, "Unable to encode notifications to JSON", http.StatusInternalServerError)
	}
}
func (handler *NotificationHandler) SetNotificationStatusDiscarded(w http.ResponseWriter, r *http.Request) {
	// Extract the user_id and notification ID from the URL path parameters
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	notificationID, ok := vars["id"]
	if !ok {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	// Parse the createdAt timestamp from the query parameters
	createdAt := r.URL.Query().Get("created_at")
	if createdAt == "" {
		http.Error(w, "CreatedAt timestamp is required", http.StatusBadRequest)
		return
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		http.Error(w, "Invalid CreatedAt timestamp format", http.StatusBadRequest)
		return
	}

	// Call the service method to set the status
	err = handler.service.SetStatusDiscarded(userID, notificationID, parsedCreatedAt)
	if err != nil {
		http.Error(w, "Failed to set status to discarded: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification status set to discarded successfully"))
}
