package startup

import (
	"userService/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var users = []*domain.User{
	{
		ID:       getObjectId("6360ed69e504b6e93f964229"),
		Name:     "John",
		Surname:  "Doe",
		Username: "Doe22",
		Email:    "john.doe@example.com",
		Password: "password123", // This should ideally be hashed in a real application
		UserRole: "Manager",
		Enabled:  true,
	},
	{
		ID:       getObjectId("6360ed69e504b6e93f964230"),
		Name:     "Jane",
		Surname:  "Smith",
		Username: "Smiti11",
		Email:    "jane.smith@example.com",
		Password: "securepass", // Example only, use hashed passwords in practice
		UserRole: "User",
		Enabled:  true,
	},
	{
		ID:       getObjectId("6360ed69e504b6e93f964231"),
		Name:     "Bob",
		Surname:  "Brown",
		Username: "Brown99",
		Email:    "bob.brown@example.com",
		Password: "$2a$10$yGa7vVtGMiNI5fpu0atTZuF/EhPElmeU83D2VQYRvAqdmL8nd93W.", // Replace with hashed passwords in production
		UserRole: "User",
		Enabled:  false, // Example of a disabled user
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
