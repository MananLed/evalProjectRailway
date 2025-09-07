package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/internal/response"
	"github.com/MananLed/evalProjectRailway/internal/service"
	"github.com/MananLed/evalProjectRailway/internal/utils"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
	"github.com/google/uuid"
)

type TicketHandler struct {
	TicketService service.TicketServiceInterface
	UserService   service.UserServiceInterface
	TrainService  service.TrainServiceInterface
}

func NewTicketHandler(ts service.TicketServiceInterface, us service.UserServiceInterface, trs service.TrainServiceInterface) *TicketHandler {
	return &TicketHandler{
		TicketService: ts,
		UserService:   us,
		TrainService:  trs,
	}
}

func (h *TicketHandler) BookTicket(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user id not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}

	var req struct {
		TrainID     string `json:"trainid"`
		BookedSeats int    `json:"bookedseats"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid Input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid Input", 1001)
		return
	}

	trainID, err := uuid.Parse(req.TrainID)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid train ID", 1002)
		return
	}

	if req.BookedSeats > 5 {
		response.ErrorResponse(w, http.StatusBadRequest, "Cannot book more than 5 seats at a time", 1002)
		return
	}

	train, err := h.TrainService.GetTrainByID(trainID)
	if err != nil {
		response.ErrorResponse(w, http.StatusNotFound, "Train with such id does not exist", 1008)
		return
	}

	if train.TotalSeats-train.BookedSeats < req.BookedSeats {
		response.ErrorResponse(w, http.StatusNotAcceptable, "Overbooking cannot be done", 1090)
		return
	}

	ticket := model.Ticket{
		ID:          utils.GenerateUUID(),
		TrainID:     trainID,
		PassengerID: user.ID,
		TotalFare:   req.BookedSeats * train.SeatFare,
		BookedSeats: req.BookedSeats,
		Source:      train.Source,
		Destination: train.Destination,
		Departure:   train.Departure,
	}

	err = h.TicketService.BookTicket(ticket)

	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to book ticket", 1010)
		return
	}

	response.SuccessResponse(w, nil, "Ticket booked successfully", http.StatusOK)
}

func (h *TicketHandler) CancelTicket(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/tickets/")
	if id == "" {
		logger.LogToFile("Ticket ID missing in URL")
		response.ErrorResponse(w, http.StatusBadRequest, "Ticket ID is required", 1001)
		return
	}

	ticketid, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid ticket ID", 1002)
		return
	}

	ticket, err := h.TicketService.GetTicketByID(ticketid)
	if err != nil {
		logger.LogToFile("Ticket not found")
		response.ErrorResponse(w, http.StatusNotFound, "Ticket not found", 1008)
		return
	}

	if ticket.PassengerID != user.ID {
		logger.LogToFile("unauthorized access")
		response.ErrorResponse(w, http.StatusForbidden, "not authorized to perform the action", 1008)
		return
	}

	journeyDate := ticket.Departure.Truncate(24 * time.Hour)
	today := time.Now().UTC().Truncate(24 * time.Hour)

	if today.Equal(journeyDate) {
		response.ErrorResponse(w, http.StatusBadRequest, "Cannot cancel ticket on the day of journey", 1001)
		return
	}

	if today.After(journeyDate) {
		response.ErrorResponse(w, http.StatusBadRequest, "Cannot cancel ticket after the journey", 1001)
		return
	}

	err = h.TicketService.CancelTicket(ticketid)

	if err != nil {
		logger.LogToFile("Error cancelling ticket")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error cancelling the ticket", 1010)
		return
	}

	response.SuccessResponse(w, nil, "Ticket cancelled successfully", http.StatusOK)
}

func (h *TicketHandler) GetTicketOfPassenger(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/tickets/")

	var tickets []model.Ticket

	if id == "" && user.Role == model.RolePassenger {
		tickets, err = h.TicketService.GetTicketsOfPassenger(user.ID)
		if err != nil {
			logger.LogToFile("Error retrieving tickets")
			response.ErrorResponse(w, http.StatusInternalServerError, "Error in retrieval of tickets", 1010)
			return
		}
	} else if id != "" && user.Role == model.RoleAdmin {
		passengerid, err := uuid.Parse(id)
		if err != nil {
			logger.LogToFile("Error retrieving tickets")
			response.ErrorResponse(w, http.StatusInternalServerError, "Error in retrieval of tickets", 1010)
			return
		}

		tickets, err = h.TicketService.GetTicketsOfPassenger(passengerid)
		if err != nil {
			logger.LogToFile("Error retrieving tickets")
			response.ErrorResponse(w, http.StatusInternalServerError, "Error in retrieval of tickets", 1010)
			return
		}
	}else {
		logger.LogToFile("Invalid operation")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid operation", 1002)
		return
	}

	response.SuccessResponse(w, tickets, "Ticket details fetched successfully", http.StatusOK)
}
