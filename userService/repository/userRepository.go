package repository

import (
	"context"
	"userService/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "usersdb" // Database name
	COLLECTION = "users"   // Collection name
)

type UserMongoDBStore struct {
	users *mongo.Collection
}

func NewUserMongoDBStore(client *mongo.Client) *UserMongoDBStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &UserMongoDBStore{
		users: users,
	}
}

func (store *UserMongoDBStore) Get(id primitive.ObjectID) (*domain.User, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetAll() ([]*domain.User, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) Insert(user *domain.User) error {
	result, err := store.users.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *UserMongoDBStore) Update(id primitive.ObjectID, updateData map[string]interface{}) error {
	update := bson.M{
		"$set": updateData,
	}
	_, err := store.users.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (store *UserMongoDBStore) Delete(id primitive.ObjectID) error {
	_, err := store.users.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (store *UserMongoDBStore) DeleteAll() {
	store.users.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *UserMongoDBStore) FindByEmail(email string) (*domain.User, error) {
	filter := bson.M{"email": email}
	return store.filterOne(filter)
}

// Internal helper for filtering multiple users
func (store *UserMongoDBStore) filter(filter interface{}) ([]*domain.User, error) {
	cursor, err := store.users.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	return decode(cursor)
}

// Internal helper for filtering a single user
func (store *UserMongoDBStore) filterOne(filter interface{}) (*domain.User, error) {
	result := store.users.FindOne(context.TODO(), filter)
	var user domain.User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Decode cursor into a slice of User objects
func decode(cursor *mongo.Cursor) ([]*domain.User, error) {
	var users []*domain.User
	for cursor.Next(context.TODO()) {
		var user domain.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	err := cursor.Err()
	return users, err
}
