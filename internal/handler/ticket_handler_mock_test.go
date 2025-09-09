package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MananLed/evalProjectRailway/internal/mocks"
	"github.com/MananLed/evalProjectRailway/internal/model"
	"go.uber.org/mock/gomock"
	"github.com/google/uuid"
)

func TestBookTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	trainSvc := mocks.NewMockTrainServiceInterface(ctrl)
	ticketSvc := mocks.NewMockTicketServiceInterface(ctrl)

	handler := NewTicketHandler(ticketSvc, userSvc, trainSvc)

	adminID := uuid.New()
	passenger := &model.User{
		ID:   adminID,
		Role: model.RolePassenger,
	}

	trainID := uuid.New()
	train := &model.Train{
		ID:          trainID,
		Source:      "Delhi",
		Destination: "Mumbai",
		TotalSeats:  100,
		BookedSeats: 50,
		SeatFare:    500,
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(passenger, nil).AnyTimes()
	trainSvc.EXPECT().GetTrainByID(trainID).Return(train, nil).AnyTimes()
	ticketSvc.EXPECT().BookTicket(gomock.Any()).Return(nil).AnyTimes()

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer([]byte(`invalid json`)))
		w := httptest.NewRecorder()
		handler.BookTicket(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Overbooking", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"trainid":     trainID.String(),
			"bookedseats": 51,
		})
		req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.BookTicket(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Successful booking", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"trainid":     trainID.String(),
			"bookedseats": 2,
		})
		req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.BookTicket(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestCancelTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	ticketSvc := mocks.NewMockTicketServiceInterface(ctrl)

	handler := NewTicketHandler(ticketSvc, userSvc, nil)

	passengerID := uuid.New()
	passenger := &model.User{
		ID:   passengerID,
		Role: model.RolePassenger,
	}

	today := time.Now().UTC()
	tomorrow := today.Add(24 * time.Hour)

	ticketID := uuid.New()
	ticket := &model.Ticket{
		ID:          ticketID,
		PassengerID: passengerID,
		Departure:   tomorrow,
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(passenger, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketByID(ticketID).Return(ticket, nil).AnyTimes()
	ticketSvc.EXPECT().CancelTicket(ticketID).Return(nil).AnyTimes()

	t.Run("Missing ticket ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/tickets/", nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Successful cancellation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/tickets/"+ticketID.String(), nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestGetTicketOfPassenger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	ticketSvc := mocks.NewMockTicketServiceInterface(ctrl)

	handler := NewTicketHandler(ticketSvc, userSvc, nil)

	passengerID := uuid.New()
	passenger := &model.User{
		ID:   passengerID,
		Role: model.RolePassenger,
	}

	tickets := []model.Ticket{
		{
			ID:          uuid.New(),
			PassengerID: passengerID,
			Source:      "Delhi",
			Destination: "Mumbai",
		},
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(passenger, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketsOfPassenger(passengerID).Return(tickets, nil).AnyTimes()

	t.Run("Get tickets for passenger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tickets/", nil)
		w := httptest.NewRecorder()
		handler.GetTicketOfPassenger(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}
