package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
)

func TestBookTicket_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	ticket := &model.Ticket{
		ID:          uuid.New(),
		TrainID:     uuid.New(),
		PassengerID: uuid.New(),
		Source:      "Delhi",
		Destination: "Mumbai",
		Departure:   time.Now(),
		TotalFare:   1000,
		BookedSeats: 2,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO tickets`).
		WithArgs(ticket.ID, ticket.TrainID, ticket.PassengerID, ticket.Source, ticket.Destination,
			ticket.Departure, ticket.TotalFare, ticket.BookedSeats).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE trains SET bookedseats`).
		WithArgs(ticket.BookedSeats, ticket.TrainID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.BookTicket(ticket)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestBookTicket_BeginError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	ticket := &model.Ticket{ID: uuid.New()}

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	err := repo.BookTicket(ticket)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestBookTicket_InsertError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	ticket := &model.Ticket{ID: uuid.New(), TrainID: uuid.New(), PassengerID: uuid.New(), Departure: time.Now()}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO tickets`).
		WithArgs(ticket.ID, ticket.TrainID, ticket.PassengerID, ticket.Source, ticket.Destination,
			ticket.Departure, ticket.TotalFare, ticket.BookedSeats).
		WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	err := repo.BookTicket(ticket)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestBookTicket_UpdateError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	ticket := &model.Ticket{ID: uuid.New(), TrainID: uuid.New(), PassengerID: uuid.New(), Departure: time.Now()}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO tickets`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE trains SET bookedseats`).
		WithArgs(ticket.BookedSeats, ticket.TrainID).
		WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	err := repo.BookTicket(ticket)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestBookTicket_CommitError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	ticket := &model.Ticket{ID: uuid.New(), TrainID: uuid.New(), PassengerID: uuid.New(), Departure: time.Now()}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO tickets`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE trains SET bookedseats`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sql.ErrTxDone)

	err := repo.BookTicket(ticket)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetTicketByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "trainid", "passengerid", "source", "destination", "departure", "totalfare", "bookedseats"}).
		AddRow(id, uuid.New(), uuid.New(), "Delhi", "Mumbai", time.Now(), 500, 1)

	mock.ExpectQuery(`SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(rows)

	ticket, err := repo.GetTicketByID(id)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ticket == nil {
		t.Fatalf("expected ticket, got nil")
	}
}

func TestGetTicketByID_NoRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()

	mock.ExpectQuery(`SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	ticket, err := repo.GetTicketByID(id)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if ticket != nil {
		t.Errorf("expected nil ticket, got %+v", ticket)
	}
}

func TestGetTicketByID_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()

	mock.ExpectQuery(`SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	ticket, err := repo.GetTicketByID(id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if ticket != nil {
		t.Errorf("expected nil ticket, got %+v", ticket)
	}
}

func TestCancelTicket_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()
	trainID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT trainid, bookedseats FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"trainid", "bookedseats"}).AddRow(trainID, 2))
	mock.ExpectExec(`DELETE FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE trains SET bookedseats = bookedseats -`).
		WithArgs(2, trainID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CancelTicket(id)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCancelTicket_BeginError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	err := repo.CancelTicket(id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCancelTicket_SelectError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT trainid, bookedseats FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	err := repo.CancelTicket(id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCancelTicket_DeleteError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()
	trainID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT trainid, bookedseats FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"trainid", "bookedseats"}).AddRow(trainID, 2))
	mock.ExpectExec(`DELETE FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	err := repo.CancelTicket(id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCancelTicket_UpdateError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	id := uuid.New()
	trainID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT trainid, bookedseats FROM tickets WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"trainid", "bookedseats"}).AddRow(trainID, 2))
	mock.ExpectExec(`DELETE FROM tickets WHERE id = \$1`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE trains SET bookedseats = bookedseats -`).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	err := repo.CancelTicket(id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetTicketsOfPassenger_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	pid := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "trainid", "passengerid", "source", "destination", "departure", "totalfare", "bookedseats"}).
		AddRow(uuid.New(), uuid.New(), pid, "Delhi", "Mumbai", time.Now(), 500, 1).
		AddRow(uuid.New(), uuid.New(), pid, "Delhi", "Pune", time.Now(), 700, 2)

	mock.ExpectQuery(`SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE passengerid = \$1`).
		WithArgs(pid).
		WillReturnRows(rows)

	tickets, err := repo.GetTicketsOfPassenger(pid)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tickets) != 2 {
		t.Errorf("expected 2 tickets, got %d", len(tickets))
	}
}

func TestGetTicketsOfPassenger_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	pid := uuid.New()

	mock.ExpectQuery(`SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE passengerid = \$1`).
		WithArgs(pid).
		WillReturnError(sql.ErrConnDone)

	tickets, err := repo.GetTicketsOfPassenger(pid)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if tickets != nil {
		t.Errorf("expected nil tickets, got %+v", tickets)
	}
}

func TestGetTicketsOfPassenger_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewTicketRepository(db)
	pid := uuid.New()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(uuid.New())

	mock.ExpectQuery(`SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE passengerid = \$1`).
		WithArgs(pid).
		WillReturnRows(rows)

	tickets, err := repo.GetTicketsOfPassenger(pid)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if tickets != nil {
		t.Errorf("expected nil tickets, got %+v", tickets)
	}
}
