package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MananLed/evalProjectRailway/internal/mocks"
	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
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

func TestBookTicketErrors(t *testing.T){
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

	trainID1 := uuid.New() 
	trainID2 := uuid.New()

	train2 := &model.Train{
		ID:          trainID2,
		Source:      "Delhi",
		Destination: "Mumbai",
		TotalSeats:  3,
		BookedSeats: 0,
		SeatFare:    500,
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(passenger, nil).AnyTimes()
	trainSvc.EXPECT().GetTrainByID(trainID).Return(train, nil).AnyTimes()
	ticketSvc.EXPECT().BookTicket(gomock.Any()).Return(errors.New("slfj")).AnyTimes()
	trainSvc.EXPECT().GetTrainByID(trainID1).Return(nil, errors.New("fldsf")).AnyTimes()
	trainSvc.EXPECT().GetTrainByID(trainID2).Return(train2, nil).AnyTimes()

	t.Run("Unsuccessful booking", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"trainid":     trainID2.String(),
			"bookedseats": 5,
		})
		req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.BookTicket(w, req)

		if w.Code != http.StatusNotAcceptable {
			t.Errorf("expected status %d, got %d", http.StatusNotAcceptable, w.Code)
		}
	})

	t.Run("Train not present", func(t *testing.T){
		body, _ := json.Marshal(map[string]interface{}{
			"trainid":     trainID1.String(),
			"bookedseats": 5,
		})
		req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.BookTicket(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("Error in ticket booking", func(t *testing.T){
		body, _ := json.Marshal(map[string]interface{}{
			"trainid" : trainID2.String(),
			"bookedseats" : 1,
		})
		req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handler.BookTicket(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
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

	ticketID2 := uuid.New()

	ticket2 := &model.Ticket{
		ID:          ticketID2,
		PassengerID: passengerID,
		Departure:   tomorrow,
	}
	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(passenger, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketByID(ticketID).Return(ticket, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketByID(ticketID2).Return(ticket2, nil).AnyTimes()
	ticketSvc.EXPECT().CancelTicket(ticketID).Return(nil).AnyTimes()
	ticketSvc.EXPECT().CancelTicket(ticketID2).Return(errors.New("dslfj")).AnyTimes()

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

	t.Run("Failure in cancellation", func(t *testing.T){
		req := httptest.NewRequest(http.MethodDelete, "/tickets/"+ticketID2.String(), nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestCancelTicketError(t *testing.T) {
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
	yesterday := today.Add(-24 * time.Hour)

	ticketID := uuid.New()
	ticket := &model.Ticket{
		ID:          ticketID,
		PassengerID: passengerID,
		Departure:   tomorrow,
	}

	ticketID2 := uuid.New()
	ticket2 := &model.Ticket{
		ID: ticketID2,
		PassengerID: uuid.New(),
		Departure: tomorrow,
	}

	ticketID3 := uuid.New()
	ticket3 := &model.Ticket{
		ID: ticketID3,
		PassengerID: passengerID,
		Departure: today,
	}

	ticketID4 := uuid.New()
	ticket4 := &model.Ticket{
		ID: ticketID4,
		PassengerID: passengerID,
		Departure: yesterday,
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(passenger, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketByID(ticketID).Return(ticket, errors.New("lsfjl")).AnyTimes()
	ticketSvc.EXPECT().CancelTicket(ticketID).Return(nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketByID(ticketID2).Return(ticket2, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketByID(ticketID3).Return(ticket3, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketByID(ticketID4).Return(ticket4, nil).AnyTimes()
	t.Run("Error from get ticket", func(t *testing.T){
		req := httptest.NewRequest(http.MethodDelete, "/tickets/"+ticketID.String(), nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("Error from mismatch id", func(t *testing.T){
		req := httptest.NewRequest(http.MethodDelete, "/tickets/"+ticketID2.String(), nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("Error because cancelled on the day", func(t *testing.T){
		req := httptest.NewRequest(http.MethodDelete, "/tickets/"+ticketID3.String(), nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
	t.Run("Error because after the day", func(t *testing.T){
		req := httptest.NewRequest(http.MethodDelete, "/tickets/"+ticketID4.String(), nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
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

func TestGetTicketOfPassengerErrorONID(t *testing.T) {
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

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(passenger, errors.New("dlfjss")).AnyTimes()
	ticketSvc.EXPECT().GetTicketsOfPassenger(passengerID).Return(tickets, nil).AnyTimes()

	t.Run("Get tickets for passenger error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tickets/", nil)
		w := httptest.NewRecorder()
		handler.GetTicketOfPassenger(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Cancel ticket passenger error", func(t *testing.T){
		req := httptest.NewRequest(http.MethodDelete, "/tickets/" + tickets[0].ID.String(), nil)
		w := httptest.NewRecorder()
		handler.CancelTicket(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestGetTicketOfPassengerErrorOnGetTicket(t *testing.T) {
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
	ticketSvc.EXPECT().GetTicketsOfPassenger(passengerID).Return(tickets, errors.New("dlfsj")).AnyTimes()

	t.Run("Get tickets for passenger error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tickets/", nil)
		w := httptest.NewRecorder()
		handler.GetTicketOfPassenger(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestGetTicketOfPassengerInvalidOperation(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	ticketSvc := mocks.NewMockTicketServiceInterface(ctrl)

	handler := NewTicketHandler(ticketSvc, userSvc, nil)

	adminID := uuid.New()
	admin := &model.User{
		ID:   adminID,
		Role: model.RoleAdmin,
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketsOfPassenger(adminID).Return(nil, nil).AnyTimes()
	
	t.Run("Invalid Operation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tickets/", nil)
		w := httptest.NewRecorder()
		handler.GetTicketOfPassenger(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestGetTicketOfPassengerByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	ticketSvc := mocks.NewMockTicketServiceInterface(ctrl)

	handler := NewTicketHandler(ticketSvc, userSvc, nil)

	passengerID := uuid.New()

	adminID := uuid.New()
	admin := &model.User{
		ID:   adminID,
		Role: model.RoleAdmin,
	}

	tickets := []model.Ticket{
		{
			ID:          uuid.New(),
			PassengerID: passengerID,
			Source:      "Delhi",
			Destination: "Mumbai",
		},
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketsOfPassenger(passengerID).Return(tickets, nil).AnyTimes()

	t.Run("Get tickets for passenger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tickets/" + passengerID.String(), nil)
		w := httptest.NewRecorder()
		handler.GetTicketOfPassenger(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestGetTicketOfPassengerByAdminGetTicketError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := mocks.NewMockUserServiceInterface(ctrl)
	ticketSvc := mocks.NewMockTicketServiceInterface(ctrl)

	handler := NewTicketHandler(ticketSvc, userSvc, nil)

	passengerID := uuid.New()

	adminID := uuid.New()
	admin := &model.User{
		ID:   adminID,
		Role: model.RoleAdmin,
	}

	tickets := []model.Ticket{
		{
			ID:          uuid.New(),
			PassengerID: passengerID,
			Source:      "Delhi",
			Destination: "Mumbai",
		},
	}

	userSvc.EXPECT().GetUserByID(gomock.Any()).Return(admin, nil).AnyTimes()
	ticketSvc.EXPECT().GetTicketsOfPassenger(passengerID).Return(tickets, errors.New("dslfj")).AnyTimes()

	t.Run("Get tickets for passenger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tickets/" + passengerID.String(), nil)
		w := httptest.NewRecorder()
		handler.GetTicketOfPassenger(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}