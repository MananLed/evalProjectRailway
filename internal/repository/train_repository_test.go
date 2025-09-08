package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
)

func TestAddTrain(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}

	defer db.Close()

	repo := NewTrainRepository(db)

	train := &model.Train{
		ID:          uuid.New(),
		Name:        "Express",
		Source:      "Delhi",
		Destination: "Mumbai",
		Departure:   time.Now(),
		Arrival:     time.Now().Add(5 * time.Hour),
		TotalSeats:  100,
		BookedSeats: 10,
		SeatFare:    500,
	}

	mock.ExpectExec(`INSERT INTO trains`).WithArgs(train.ID, train.Name, train.Source, train.Destination, train.Departure, train.Arrival, train.TotalSeats, train.BookedSeats, train.SeatFare).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddTrain(train)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteTrain(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewTrainRepository(db)

	id := uuid.New()

	mock.ExpectExec(`DELETE FROM trains WHERE id = \$1`).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteTrain(id)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetTrains(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewTrainRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "name", "source", "destination", "departure",
		"arrival", "totalseats", "bookedseats", "seatfare",
	}).AddRow(uuid.New(), "Express", "Delhi", "Mumbai",
		time.Now(), time.Now().Add(5*time.Hour), 100, 10, 500)

	mock.ExpectQuery(`SELECT id, name, source, destination, departure, arrival, totalseats, bookedseats, seatfare FROM trains`).
		WillReturnRows(rows)

	trains, err := repo.GetTrains(time.Time{}, "", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(trains) != 1 {
		t.Errorf("expected 1 train, got %d", len(trains))
	}
	if trains[0].Name != "Express" {
		t.Errorf("expected name 'Express', got %s", trains[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetTrainByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewTrainRepository(db)
	id := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id", "name", "source", "destination", "departure",
		"arrival", "totalseats", "bookedseats", "seatfare",
	}).AddRow(id, "Express", "Delhi", "Mumbai",
		time.Now(), time.Now().Add(5*time.Hour), 100, 10, 500)

	mock.ExpectQuery(`SELECT id, name, source, destination, departure, arrival, totalseats, bookedseats, seatfare FROM trains WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(rows)

	train, err := repo.GetTrainByID(id)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if train == nil {
		t.Fatalf("expected train, got nil")
	}
	if train.Name != "Express" {
		t.Errorf("expected name 'Express', got %s", train.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetTrainByID_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewTrainRepository(db)
	id := uuid.New()

	mock.ExpectQuery(`SELECT id, name, source, destination, departure, arrival, totalseats, bookedseats, seatfare FROM trains WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	train, err := repo.GetTrainByID(id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if train != nil {
		t.Errorf("expected nil train, got %+v", train)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
