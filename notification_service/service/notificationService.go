package service

import (
	"crypto/md5"
	"fmt"
	"notificationService/domain"
	"notificationService/repository"
	"time"

	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService struct {
	repo *repository.NotificationCassandraStore
}

// Constructor for NotificationService
func NewNotificationService(repo *repository.NotificationCassandraStore) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (service *NotificationService) GetAll() ([]*domain.Notification, error) {
	// Call the repository method to fetch all notifications
	return service.repo.GetAll()
}

// Fetch a single notification by ID
func (service *NotificationService) Get(id string) (*domain.Notification, error) {
	// Convert the string ID to gocql.UUID
	notificationID, err := gocql.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	// Fetch the notification from the repository
	return service.repo.Get(notificationID)
}

func objectIDToUUID(objectID primitive.ObjectID) gocql.UUID {
	hash := md5.Sum(objectID[:])
	uid, _ := gocql.ParseUUID(fmt.Sprintf("%x-%x-%x-%x-%x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16]))
	return uid
}

func (service *NotificationService) GetByUserId(userId string) ([]*domain.Notification, error) {
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	userUUID := objectIDToUUID(objectID)

	return service.repo.GetByUserId(userUUID)
}
func (service *NotificationService) SetStatusDiscarded(userID string, notificationID string, createdAt time.Time) error {
	// Convert the string notificationID to gocql.UUID
	id, err := gocql.ParseUUID(notificationID)
	if err != nil {
		return err
	}

	// Convert the userID to gocql.UUID
	userIDUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return err
	}

	// Call the repository method to update the status
	return service.repo.SetStatusDiscarded(userIDUUID, createdAt, id)
}

// Delete all notifications (use cautiously in production environments)
func (service *NotificationService) DeleteAll() error {
	return service.repo.DeleteAll()
}
