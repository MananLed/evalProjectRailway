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

func (r *TrainRepository) DeleteTrain(id uuid.UUID) error{
	query := `DELETE FROM trains WHERE id = $1`

	r.mu.Lock()
	_, err := r.db.Exec(query, id)
	r.mu.Unlock()

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	return nil 
}

func (r *TrainRepository) GetTrains(travelDate time.Time, source string, destinaiton string) ([]model.Train, error){
	query := `SELECT id, name, source, destination, departure, arrival, totalseats, bookedseats, seatfare FROM trains WHERE 1=1`
	args := []interface{}{}

	argCount := 1

	if !travelDate.IsZero() {
		query += fmt.Sprintf(" AND DATE(departure) = $%d", argCount)
		args = append(args, travelDate.Format("2006-01-02"))
		argCount++
	}

	if source != ""{
		query += fmt.Sprintf(" AND source = $%d", argCount)
		args = append(args, source)
		argCount++
	}

	if destinaiton != "" {
		query += fmt.Sprintf(" AND destination = $%d", argCount)
		args = append(args, destinaiton)
		argCount++
	}

	r.mu.Lock()
	rows, err := r.db.Query(query, args...)
	r.mu.Unlock()

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err 
	}
	defer rows.Close()

	var trains []model.Train

	for rows.Next() {
		var train model.Train

		err := rows.Scan(&train.ID, &train.Name, &train.Source, &train.Destination, &train.Departure, &train.Arrival, &train.TotalSeats, &train.BookedSeats, &train.SeatFare)

		if err != nil{
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, err 
		}
		trains = append(trains, train)
	}
	return trains, nil 
}

func (r *TrainRepository) GetTrainByID(id uuid.UUID) (*model.Train, error) {
	query := `SELECT id, name, source, destination, departure, arrival, totalseats, bookedseats, seatfare FROM trains WHERE id = $1`

	var train model.Train

	r.mu.Lock()
	err := r.db.QueryRow(query, id).Scan(&train.ID, &train.Name, &train.Source, &train.Destination, &train.Departure, &train.Arrival, &train.TotalSeats, &train.BookedSeats, &train.SeatFare)
	r.mu.Unlock()

	if err != nil{
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("error: %v", err)
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	return &train, nil 
}

