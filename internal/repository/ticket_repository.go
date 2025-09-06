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
	query := `
		INSERT INTO tickets (id, trainid, passengerid, source, destination, departure, totalfare, bookedseats)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	r.mu.Lock()
	_, err := r.db.Exec(query, ticket.ID, ticket.TrainID, ticket.PassengerID, ticket.Source, ticket.Destination, ticket.Departure, ticket.TotalFare, ticket.BookedSeats)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return nil
}
