package repository

import (
	"context"
	"fmt"
	"log"
	"taskService/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepo struct {
	client  *mongo.Client
	logger  *log.Logger
	dbName  string
	colName string
}

func (tr *TaskRepo) getCollection() *mongo.Collection {
	return tr.client.Database(tr.dbName).Collection(tr.colName)
}

func NewTaskRepo(ctx context.Context, logger *log.Logger) (*TaskRepo, error) {
	dburi := "your_mongo_connection_string" // Replace with your MongoDB URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	return &TaskRepo{
		client:  client,
		logger:  logger,
		dbName:  "tasksdb",
		colName: "tasks",
	}, nil
}
func (pr *TaskRepo) Disconnect(ctx context.Context) error {
	err := pr.client.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Create adds a new task to the database
func (tr *TaskRepo) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	collection := tr.getCollection()

	// Assigning a new ObjectID if not already set
	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}

	// Insert the task into the collection
	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}

	return task, nil
}

// Update updates an existing task in the database
func (tr *TaskRepo) Update(ctx context.Context, taskID string, updatedTask *domain.Task) (*domain.Task, error) {
	collection := tr.getCollection()

	// Convert string taskID to ObjectID
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %v", err)
	}

	// Create filter and update for the task
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        updatedTask.Name,
			"description": updatedTask.Description,
			"status":      updatedTask.Status,
			"user":        updatedTask.UserID,
		},
	}

	// Perform the update
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %v", err)
	}

	// Check if task was updated
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("task not found")
	}

	// Return the updated task
	updatedTask.ID = id
	return updatedTask, nil
}

// Delete deletes a task from the database
func (tr *TaskRepo) Delete(ctx context.Context, taskID string) error {
	collection := tr.getCollection()

	// Convert string taskID to ObjectID
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	// Perform the deletion
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	// Check if any task was deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
