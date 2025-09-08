package service

import (
	"testing"
	"time"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
)

type MockTrainRepo struct {
	trains map[uuid.UUID]model.Train
}

func (m *MockTrainRepo) AddTrain(train *model.Train) error {
	return nil
}

func (m *MockTrainRepo) DeleteTrain(id uuid.UUID) error {
	return nil
}

func (m *MockTrainRepo) GetTrains(travelDate time.Time, source string, destination string) ([]model.Train, error) {
	var mt []model.Train
	return mt, nil
}

func (m *MockTrainRepo) GetTrainByID(id uuid.UUID) (*model.Train, error) {
	var mt model.Train

	return &mt, nil
}

func Test_AddTrain(t *testing.T) {
	mockRepo := &MockTrainRepo{trains: make(map[uuid.UUID]model.Train)}

	service := NewTrainService(mockRepo)

	var train model.Train
	err := service.AddTrain(train) 

	if err != nil{
		t.Error("Expected success, but got error in Add Train")
	}
}
