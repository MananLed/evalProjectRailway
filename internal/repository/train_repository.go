package repository

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
	"github.com/google/uuid"
)

type TrainRepositoryInterface interface {
	AddTrain(train *model.Train) error
	DeleteTrain(id uuid.UUID) error 
	GetTrains(travelDate time.Time, source string, destination string) ([]model.Train, error)
	GetTrainByID(id uuid.UUID) (*model.Train, error)
}

type TrainRepository struct {
	mu sync.Mutex
	db *sql.DB
}

func NewTrainRepository(db *sql.DB) *TrainRepository {
	return &TrainRepository{
		db: db,
	}
}

func (r *TrainRepository) AddTrain(train *model.Train) error {

	query := `
		INSERT INTO trains (id, name, source, destination, departure, arrival, totalseats, bookedseats, seatfare)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	r.mu.Lock()
	_, err := r.db.Exec(query, train.ID, train.Name, train.Source, train.Destination, train.Departure, train.Arrival, train.TotalSeats, train.BookedSeats, train.SeatFare)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return nil
}
