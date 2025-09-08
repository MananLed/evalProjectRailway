//go:generate mockgen -source=user_service.go -destination=../mocks/user_mock_service.go -package=mocks
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/internal/repository"
	"github.com/MananLed/evalProjectRailway/internal/utils"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	IsEmailUnique(email string) bool
	SignUp(user model.User) error
	Login(email string, password string) (*model.User, error)
	GetUserByID(ctx context.Context) (*model.User, error)
	UpdateProfile(user model.User) error
	ChangePassword(user *model.User, currentPassword string, newPassword string) error
	DeleteProfile(id uuid.UUID) error
}

type UserService struct {
	UserRepository repository.UserRepositoryInterface
}

func NewUserService(ur repository.UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: ur,
	}
}

func (s *UserService) IsEmailUnique(email string) bool {
	return s.UserRepository.IsEmailUnique(email)
}

func (s *UserService) SignUp(user model.User) error {
	return s.UserRepository.AddUser(user)
}

func (s *UserService) Login(email string, password string) (*model.User, error) {
	user, err := s.UserRepository.GetUserByEmailAndPassword(email, password)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context) (*model.User, error) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	return s.UserRepository.GetUserByID(user.ID)
}

func (s *UserService) UpdateProfile(user model.User) error {
	return s.UserRepository.UpdateUser(user)
}

func (s *UserService) ChangePassword(user *model.User, currentPassword string, newPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return errors.New("current password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		return errors.New("new password has no change")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return s.UserRepository.ChangePassword(user.ID, string(hashedPassword))
}

func (s *UserService) DeleteProfile(id uuid.UUID) error {
	return s.UserRepository.DeleteUserByID(id)
}
