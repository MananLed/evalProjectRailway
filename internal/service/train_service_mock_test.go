package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
)

type MockTrainRepo struct {
	trains map[uuid.UUID]model.Train
}

func (m *MockTrainRepo) AddTrain(train *model.Train) error {
	m.trains[train.ID] = *train
	return nil
}

func (m *MockTrainRepo) DeleteTrain(id uuid.UUID) error {
	if _, exists := m.trains[id]; !exists {
		return fmt.Errorf("train not found")
	}
	delete(m.trains, id)
	return nil
}

func (m *MockTrainRepo) GetTrains(travelDate time.Time, source string, destination string) ([]model.Train, error) {
	var result []model.Train
	for _, t := range m.trains {
		if t.Source == source && t.Destination == destination && t.Departure.Truncate(24*time.Hour).Equal(travelDate.Truncate(24*time.Hour)) {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *MockTrainRepo) GetTrainByID(id uuid.UUID) (*model.Train, error) {
	train, exists := m.trains[id]
	if !exists {
		return nil, fmt.Errorf("train not found")
	}
	return &train, nil
}

func setupService() (*TrainService, *MockTrainRepo) {
	mockRepo := &MockTrainRepo{trains: make(map[uuid.UUID]model.Train)}
	service := NewTrainService(mockRepo)
	return service, mockRepo
}

func TestAddTrain(t *testing.T) {
	svc, repo := setupService()
	train := model.Train{
		ID:          uuid.New(),
		Name:        "Express",
		Source:      "Delhi",
		Destination: "Mumbai",
		Departure:   time.Now(),
	}

	err := svc.AddTrain(train)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	stored, exists := repo.trains[train.ID]
	if !exists {
		t.Errorf("expected train to be stored, but not found")
	}
	if stored.Name != "Express" {
		t.Errorf("expected train name 'Express', got '%s'", stored.Name)
	}
}

func TestDeleteTrain(t *testing.T) {
	svc, repo := setupService()
	id := uuid.New()

	repo.trains[id] = model.Train{ID: id}
	err := svc.DeleteTrain(id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if _, exists := repo.trains[id]; exists {
		t.Errorf("expected train to be deleted, but still exists")
	}

	err = svc.DeleteTrain(uuid.New())
	if err == nil {
		t.Errorf("expected error for deleting non-existing train, got nil")
	}
}

func TestGetTrains(t *testing.T) {
	svc, repo := setupService()
	date := time.Now().Truncate(24 * time.Hour)

	train1 := model.Train{ID: uuid.New(), Source: "Delhi", Destination: "Mumbai", Departure: date}
	train2 := model.Train{ID: uuid.New(), Source: "Delhi", Destination: "Mumbai", Departure: date}
	train3 := model.Train{ID: uuid.New(), Source: "Delhi", Destination: "Chennai", Departure: date}

	repo.trains[train1.ID] = train1
	repo.trains[train2.ID] = train2
	repo.trains[train3.ID] = train3

	trains, err := svc.GetTrains(date, "Delhi", "Mumbai")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(trains) != 2 {
		t.Errorf("expected 2 trains, got %d", len(trains))
	}

	trains, err = svc.GetTrains(date, "Delhi", "Kolkata")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(trains) != 0 {
		t.Errorf("expected 0 trains, got %d", len(trains))
	}
}

func TestGetTrainByID(t *testing.T) {
	svc, repo := setupService()
	id := uuid.New()

	repo.trains[id] = model.Train{ID: id, Name: "Rajdhani"}
	train, err := svc.GetTrainByID(id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if train.Name != "Rajdhani" {
		t.Errorf("expected train name 'Rajdhani', got '%s'", train.Name)
	}

	_, err = svc.GetTrainByID(uuid.New())
	if err == nil {
		t.Errorf("expected error for non-existing train, got nil")
	}
}
