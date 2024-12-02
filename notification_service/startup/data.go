package startup

import (
	"crypto/md5"
	"fmt"
	"notificationService/domain"
	"time"

	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var notifications = []*domain.Notification{
	{
		Message:   "Your order has been shipped.",
		UserID:    getUUID("674cb80b21ff024da4e0953e"),
		Status:    "unread",
		CreatedAt: time.Now(),
	},
	{
		Message:   "You have a new message.",
		UserID:    getUUID("674cb80b21ff024da4e0953e"),
		Status:    "unread",
		CreatedAt: time.Now(),
	},
	{
		Message:   "Your password has been changed successfully.",
		UserID:    getUUID("674cb80b21ff024da4e0953e"),
		Status:    "unread",
		CreatedAt: time.Now(),
	},
	{
		Message:   "Your account subscription has been renewed.",
		UserID:    getUUID("6360ed69e504b6e93f964230"),
		Status:    "unread",
		CreatedAt: time.Now(),
	},
	{
		Message:   "Your support ticket has been resolved.",
		UserID:    getUUID("6360ed69e504b6e93f964230"),
		Status:    "unread",
		CreatedAt: time.Now(),
	},
}

func objectIDToUUID(objectID primitive.ObjectID) gocql.UUID {
	hash := md5.Sum(objectID[:])
	uid, _ := gocql.ParseUUID(fmt.Sprintf("%x-%x-%x-%x-%x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16]))
	return uid
}

// Helper function to create UUID from a string.
func getUUID(id string) gocql.UUID {

	objectID, _ := primitive.ObjectIDFromHex(id)

	userUUID := objectIDToUUID(objectID)

	return userUUID
}
