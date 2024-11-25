package repository

import (
	"context"
	"fmt"
	"projectService/domain"

	// NoSQL: module containing Mongo api client
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "projectsdb"
	COLLECTION = "projects"
)

type ProjectMongoDBStore struct {
	projects *mongo.Collection
}

func NewProjectMongoDBStore(client *mongo.Client) *ProjectMongoDBStore {
	projects := client.Database(DATABASE).Collection(COLLECTION)
	return &ProjectMongoDBStore{
		projects: projects,
	}
}

func (store *ProjectMongoDBStore) Get(id primitive.ObjectID) (*domain.Project, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *ProjectMongoDBStore) GetAll() ([]*domain.Project, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *ProjectMongoDBStore) Insert(project *domain.Project) error {
	if project.Members == nil {
		project.Members = []*domain.User{}
	}

	result, err := store.projects.InsertOne(context.TODO(), project)
	if err != nil {
		return err
	}
	project.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *ProjectMongoDBStore) DeleteAll() {
	store.projects.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *ProjectMongoDBStore) filter(filter interface{}) ([]*domain.Project, error) {
	cursor, err := store.projects.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *ProjectMongoDBStore) filterOne(filter interface{}) (*domain.Project, error) {
	result := store.projects.FindOne(context.TODO(), filter)
	var project domain.Project
	err := result.Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func decode(cursor *mongo.Cursor) ([]*domain.Project, error) {
	var projects []*domain.Project
	for cursor.Next(context.TODO()) {
		var project domain.Project
		err := cursor.Decode(&project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}
	err := cursor.Err()
	return projects, err
}

func (store *ProjectMongoDBStore) AddUserToProject(projectID primitive.ObjectID, user *domain.User) error {
	// Check if the user is already a member of the project
	filter := bson.M{
		"_id": projectID,
		"members": bson.M{
			"$elemMatch": bson.M{
				"_id": user.ID,
			},
		},
	}

	existing := store.projects.FindOne(context.TODO(), filter)
	if existing.Err() == nil {
		return fmt.Errorf("already exists in project")
	}

	// Add the user to the project's members list
	update := bson.M{
		"$addToSet": bson.M{
			"members": user,
		},
	}

	// Apply the update
	_, err := store.projects.UpdateOne(context.TODO(), bson.M{"_id": projectID}, update)
	if err != nil {
		return err
	}
	return nil
}
