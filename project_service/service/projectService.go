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

func (service *ProjectService) AddUserToProject(projectId string, userId string, loggedUserID string) error {
	projectID, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return err
	}

	project, err := service.repo.Get(projectID)
	if err != nil {
		return err
	}

	if project == nil {
		return fmt.Errorf("project not found")
	}
	if project.ManagerID.Hex() != loggedUserID {
		return fmt.Errorf("user is not the project manager")
	}

	user, err := service.userClient.Get(userId)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	if len(project.Members) >= project.MaxMembers {
		return fmt.Errorf("cannot add user: project has reached max members")
	}

	err = service.repo.AddUserToProject(projectID, user)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProjectService) RemoveUserFromProject(projectID primitive.ObjectID, userID primitive.ObjectID) error {
	project, err := service.repo.Get(projectID)
	if err != nil {
		return fmt.Errorf("failed to fetch project: %v", err)
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}
	userExists := false
	for _, member := range project.Members {
		if member.ID == userID {
			userExists = true
			break
		}
	}
	if !userExists {
		return fmt.Errorf("user not part of the project")
	}
	err = service.repo.RemoveUserFromProject(projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from project: %v", err)
	}

	return nil
}

func (service *ProjectService) GetByUserId(userId string) ([]*domain.Project, error) {
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	projects, err := service.repo.GetByUserId(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching projects for user: %v", err)
	}

	return projects, nil
}
