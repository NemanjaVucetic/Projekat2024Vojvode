package service

import (
	"errors"
	"userService/domain"
	"userService/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo repository.UserMongoDBStore
}

func NewUserService(repo repository.UserMongoDBStore) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (service *UserService) Get(id string) (*domain.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return service.repo.Get(userID)
}

func (service *UserService) GetAll() ([]*domain.User, error) {
	return service.repo.GetAll()
}

func (service *UserService) Create(user *domain.User) error {
	existingUser, err := service.repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already in use")
	}
	return service.repo.Insert(user)
}

func (service *UserService) Delete(id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return service.repo.Delete(userID)
}

func (service *UserService) Update(id string, updateData map[string]interface{}) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return service.repo.Update(userID, updateData)
}

func (service *UserService) Login(email string, password string) (*domain.User, error) {
	user, err := service.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Password != password {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}
