package startup

import (
	"projectService/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var projects = []*domain.Project{}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
