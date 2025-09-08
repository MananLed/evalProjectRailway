package service

import (
	"errors"
	"testing"
	"time"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
)

type MockTicketRepo struct {
	tickets map[uuid.UUID]model.Ticket
}

func NewMockTicketRepo() *MockTicketRepo {
	return &MockTicketRepo{tickets: make(map[uuid.UUID]model.Ticket)}
}

func (m *MockTicketRepo) BookTicket(ticket *model.Ticket) error {

	if _, exists := m.tickets[ticket.ID]; exists {
		return errors.New("ticket already exists")
	}
	m.tickets[ticket.ID] = *ticket
	return nil
}

func (m *MockTicketRepo) GetTicketByID(id uuid.UUID) (*model.Ticket, error) {

	t, exists := m.tickets[id]
	if !exists {
		return nil, errors.New("ticket not found")
	}
	return &t, nil
}

func (m *MockTicketRepo) CancelTicket(id uuid.UUID) error {

	if _, exists := m.tickets[id]; !exists {
		return errors.New("ticket not found")
	}
	delete(m.tickets, id)
	return nil
}

func (m *MockTicketRepo) GetTicketsOfPassenger(id uuid.UUID) ([]model.Ticket, error) {

	var result []model.Ticket
	for _, t := range m.tickets {
		if t.PassengerID == id {
			result = append(result, t)
		}
	}
	return result, nil
}

func setupTicketService() (*TicketService, *MockTicketRepo) {
	repo := NewMockTicketRepo()
	service := NewTicketService(repo)
	return service, repo
}

func TestBookTicket(t *testing.T) {
	svc, _ := setupTicketService()
	ticket := model.Ticket{
		ID:          uuid.New(),
		TrainID:     uuid.New(),
		PassengerID: uuid.New(),
		Source:      "Delhi",
		Destination: "Mumbai",
		Departure:   time.Now(),
		TotalFare:   500,
		BookedSeats: 2,
	}

	err := svc.BookTicket(ticket)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = svc.BookTicket(ticket)
	if err == nil {
		t.Errorf("expected error when booking duplicate ticket, got nil")
	}
}

func TestGetTicketByID(t *testing.T) {
	svc, repo := setupTicketService()
	id := uuid.New()
	ticket := model.Ticket{ID: id, Source: "Delhi", Destination: "Pune"}
	repo.tickets[id] = ticket

	got, err := svc.GetTicketByID(id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if got.Source != "Delhi" {
		t.Errorf("expected Source=Delhi, got %s", got.Source)
	}

	_, err = svc.GetTicketByID(uuid.New())
	if err == nil {
		t.Errorf("expected error for non-existing ticket, got nil")
	}
}

func TestCancelTicket(t *testing.T) {
	svc, repo := setupTicketService()
	id := uuid.New()
	repo.tickets[id] = model.Ticket{ID: id}

	err := svc.CancelTicket(id)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if _, exists := repo.tickets[id]; exists {
		t.Errorf("expected ticket to be deleted, but still exists")
	}
	
	err = svc.CancelTicket(uuid.New())
	if err == nil {
		t.Errorf("expected error when canceling non-existing ticket, got nil")
	}
}

func TestGetTicketsOfPassenger(t *testing.T) {
	svc, repo := setupTicketService()
	passengerID := uuid.New()
	t1 := model.Ticket{ID: uuid.New(), PassengerID: passengerID, Source: "Delhi"}
	t2 := model.Ticket{ID: uuid.New(), PassengerID: passengerID, Source: "Mumbai"}
	t3 := model.Ticket{ID: uuid.New(), PassengerID: uuid.New(), Source: "Chennai"} // different passenger

	repo.tickets[t1.ID] = t1
	repo.tickets[t2.ID] = t2
	repo.tickets[t3.ID] = t3

	tickets, err := svc.GetTicketsOfPassenger(passengerID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(tickets) != 2 {
		t.Errorf("expected 2 tickets, got %d", len(tickets))
	}
}
