package startup

import (
	"userService/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var users = []*domain.User{}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
