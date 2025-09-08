//go:generate mockgen -source=train_service.go -destination=../mocks/train_mock_service.go -package=mocks
package service

import (
	"time"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/internal/repository"
	"github.com/google/uuid"
)

type TrainServiceInterface interface {
	AddTrain(train model.Train) error
	DeleteTrain(id uuid.UUID) error
	GetTrains(travelDate time.Time, source string, destination string) ([]model.Train, error)
	GetTrainByID(id uuid.UUID) (*model.Train, error)
}

type TrainService struct {
	TrainRepository repository.TrainRepositoryInterface
}

func NewTrainService(tr repository.TrainRepositoryInterface) *TrainService {
	return &TrainService{
		TrainRepository: tr,
	}
}

func (s *TrainService) AddTrain(train model.Train) error {
	return s.TrainRepository.AddTrain(&train)
}

func (s *TrainService) DeleteTrain(id uuid.UUID) error {
	return s.TrainRepository.DeleteTrain(id)
}

func (s *TrainService) GetTrains(travelDate time.Time, source string, destination string) ([]model.Train, error) {
	return s.TrainRepository.GetTrains(travelDate, source, destination)
}

func (s *TrainService) GetTrainByID(id uuid.UUID) (*model.Train, error) {
	return s.TrainRepository.GetTrainByID(id)
}
