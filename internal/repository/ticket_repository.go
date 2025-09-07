package repository

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
	"github.com/google/uuid"
)

type TicketRepositoryInterface interface {
	BookTicket(ticket *model.Ticket) error
	GetTicketByID(id uuid.UUID) (*model.Ticket, error)
	CancelTicket(id uuid.UUID) error 
	GetTicketsOfPassenger(id uuid.UUID) ([]model.Ticket, error)
}

type TicketRepository struct {
	mu sync.Mutex
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) BookTicket(ticket *model.Ticket) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err 
	}

	defer func(){
		if err != nil{
			tx.Rollback()
		}
	}()

	insertQuery := `INSERT INTO tickets (id, trainid, passengerid, source, destination, departure, totalfare, bookedseats) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	r.mu.Lock()
	_, err = tx.Exec(insertQuery, ticket.ID, ticket.TrainID, ticket.PassengerID, ticket.Source, ticket.Destination, ticket.Departure, ticket.TotalFare, ticket.BookedSeats)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	updateQuery := `UPDATE trains SET bookedseats = bookedseats + $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, ticket.BookedSeats, ticket.TrainID)

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	if err = tx.Commit(); err != nil {
		return err 
	}

	return nil 
}

func (r *TicketRepository) GetTicketByID(id uuid.UUID) (*model.Ticket, error) {
	
	query := `SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE id = $1`
	var ticket model.Ticket

	r.mu.Lock()
	err := r.db.QueryRow(query, id).Scan(&ticket.ID, &ticket.TrainID, &ticket.PassengerID, &ticket.Source, &ticket.Destination, &ticket.Departure, &ticket.TotalFare, &ticket.BookedSeats)
	r.mu.Unlock()

	if err != nil{
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err 
	}

	return &ticket, nil 
}

func (r *TicketRepository) CancelTicket(id uuid.UUID) error {
	tx, err := r.db.Begin()

	if err != nil{
		return err 
	}

	defer func(){
		if err != nil{
			tx.Rollback()
		}
	}()

	var trainID uuid.UUID

	var bookedSeats int 

	selectQuery := `SELECT trainid, bookedseats FROM tickets WHERE id = $1`
	err = tx.QueryRow(selectQuery, id).Scan(&trainID, &bookedSeats)

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	deleteQuery := `DELETE FROM tickets WHERE id = $1`

	_, err = tx.Exec(deleteQuery, id)

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	updateQuery := `UPDATE trains SET bookedseats = bookedseats - $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, bookedSeats, trainID)

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	if err = tx.Commit(); err != nil{
		return err 
	}

	return nil 
}

func (r *TicketRepository) GetTicketsOfPassenger(id uuid.UUID) ([]model.Ticket, error){

	query := `SELECT id, trainid, passengerid, source, destination, departure, totalfare, bookedseats FROM tickets WHERE passengerid = $1`

	r.mu.Lock()
	rows, err := r.db.Query(query, id)
	r.mu.Unlock()

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err 
	}

	defer rows.Close()

	var tickets []model.Ticket

	for rows.Next() {
		var ticket model.Ticket

		if err := rows.Scan(&ticket.ID, &ticket.TrainID, &ticket.PassengerID, &ticket.Source, &ticket.Destination, &ticket.Departure, &ticket.TotalFare, &ticket.BookedSeats); err != nil{
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, err 
		}

		tickets = append(tickets, ticket)
	}
	return tickets, nil 
}


