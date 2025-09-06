package service

import (
	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/internal/repository"
	"github.com/google/uuid"
)

type TicketServiceInterface interface {
	BookTicket(ticket model.Ticket) error
	GetTicketByID(id uuid.UUID) (*model.Ticket, error)
	CancelTicket(id uuid.UUID) error
	GetTicketsOfPassenger(id uuid.UUID) ([]model.Ticket, error)
}

type TicketService struct {
	TicketRepository repository.TicketRepositoryInterface
}

func NewTicketService(tr repository.TicketRepositoryInterface) *TicketService {
	return &TicketService{
		TicketRepository: tr,
	}
}

func (s *TicketService) BookTicket(ticket model.Ticket) error {
	return s.TicketRepository.BookTicket(&ticket)
}

func (s *TicketService) GetTicketByID(id uuid.UUID) (*model.Ticket, error) {
	return s.TicketRepository.GetTicketByID(id)
}

func (s *TicketService) CancelTicket(id uuid.UUID) error {
	return s.TicketRepository.CancelTicket(id)
}

func (s *TicketService) GetTicketsOfPassenger(id uuid.UUID) ([]model.Ticket, error) {
	return s.TicketRepository.GetTicketsOfPassenger(id)
}
