package service

import (
	"errors"
	"fmt"
	"userService/domain"
	"userService/repository"
	"userService/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
	// Check if email is already in use
	existingUser, err := service.repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Create user with `Enabled` set to false
	user.Enabled = false
	err = service.repo.Insert(user)
	if err != nil {
		return err
	}

	err = utils.SendEmail(user, "Account Confirmation")
	if err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return nil
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
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.Enabled {
		return nil, errors.New("not verified")
	}

	return user, nil
}

func (service *UserService) ValidateAccount(id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user, err := service.repo.Get(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if user.Enabled {
		return errors.New("account already verified")
	}

	updateData := map[string]interface{}{"enabled": true}
	return service.repo.Update(userID, updateData)
}

