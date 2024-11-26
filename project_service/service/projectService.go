package service

import (
	"fmt"
	"projectService/client"
	"projectService/domain"
	"projectService/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectService struct {
	repo       repository.ProjectMongoDBStore
	userClient client.Client
}

func NewProjectService(repo repository.ProjectMongoDBStore, userClient client.Client) *ProjectService {
	return &ProjectService{
		repo:       repo,
		userClient: userClient,
	}
}

func (service *ProjectService) Get(id string) (*domain.Project, error) {
	projectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return service.repo.Get(projectID)
}

func (service *ProjectService) GetAll() ([]*domain.Project, error) {
	return service.repo.GetAll()
}

func (service *ProjectService) Create(project *domain.Project) error {
	return service.repo.Insert(project)
}

func (service *ProjectService) AddUserToProject(projectId string, userId string) error {
	// Convert project ID to ObjectID
	projectID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return err
	}

	// Retrieve the project from the repository
	project, err := service.repo.Get(projectID)
	if err != nil {
		return err
	}

	// Ensure project is found
	if project == nil {
		return fmt.Errorf("project not found")
	}

	// Retrieve the user from the user service
	user, err := service.userClient.Get(userId)
	if err != nil {
		return err
	}

	// Ensure user is found
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if the project is full
	if len(project.Members) >= project.MaxMembers {
		return fmt.Errorf("cannot add user: project has reached max members")
	}

	// Add the user to the project using the repository
	err = service.repo.AddUserToProject(projectID, user)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProjectService) GetByUserId(userId string) ([]*domain.Project, error) {
	// Convert the userId string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	// Fetch projects using the repository method
	projects, err := service.repo.GetByUserId(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching projects for user: %v", err)
	}

	return projects, nil
}
