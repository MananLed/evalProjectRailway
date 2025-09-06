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

type TrainHandler struct {
	TrainService service.TrainServiceInterface
	UserService  service.UserServiceInterface
}

func NewTrainHandler(ts service.TrainServiceInterface, us service.UserServiceInterface) *TrainHandler {
	return &TrainHandler{
		TrainService: ts,
		UserService:  us,
	}
}

func (h *TrainHandler) AddTrain(w http.ResponseWriter, r *http.Request) {

	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil || (user.Role != model.RoleAdmin) {
		logger.LogToFile("unauthorized person wants to view list of residents")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	type TrainInput struct {
		Name          string `json:"name"`
		Source        string `json:"source"`
		Destination   string `json:"destination"`
		DepartureDate string `json:"departuredate"`
		DepartureTime string `json:"departuretime"`
		ArrivalDate   string `json:"arrivaldate"`
		ArrivalTime   string `json:"arrivaltime"`
		TotalSeats    int    `json:"totalseats"`
		SeatFare      int    `json:"seatfare"`
	}

	var input TrainInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	depDate, err := time.Parse("2006-01-02", input.DepartureDate)

	if err != nil {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	depTime, err := time.Parse("15:04", input.DepartureTime)

	if err != nil {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	arrDate, err := time.Parse("2006-01-02", input.ArrivalDate)

	if err != nil {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	arrTime, err := time.Parse("15:04", input.ArrivalTime)

	if err != nil {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	departure := time.Date(depDate.Year(), depDate.Month(), depDate.Day(), depTime.Hour(), depTime.Minute(), 0, 0, time.UTC)
	arrival := time.Date(arrDate.Year(), arrDate.Month(), arrDate.Day(), arrTime.Hour(), arrTime.Minute(), 0, 0, time.UTC)

	if !departure.Before(arrival) {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	train := model.Train{
		ID:          utils.GenerateUUID(),
		Name:        input.Name,
		Source:      input.Source,
		Destination: input.Destination,
		Departure:   departure,
		Arrival:     arrival,
		TotalSeats:  input.TotalSeats,
		BookedSeats: 0,
		SeatFare:    input.SeatFare,
	}

	err = h.TrainService.AddTrain(train)

	if err != nil {
		logger.LogToFile("Error adding train")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error adding train", 1006)
		return
	}

	logger.LogToFile("Train added successfully")
	response.SuccessResponse(w, nil, "Train added successfully", http.StatusCreated)
}

func (h *TrainHandler) DeleteTrain(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil || (user.Role != model.RoleAdmin) {
		logger.LogToFile("unauthorized person wants to view list of residents")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/trains/")

	if id == "" {
		logger.LogToFile("Train ID missing in URL")
		response.ErrorResponse(w, http.StatusBadRequest, "Train ID is required", 1001)
		return
	}

	trainID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request ID", 1002)
		return
	}

	err = h.TrainService.DeleteTrain(trainID)

	if err != nil {
		logger.LogToFile("Error deleting the train")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error deleting the train", 1010)
		return
	}

	response.SuccessResponse(w, nil, "Train deleted successfully", http.StatusOK)
}

func (h *TrainHandler) ViewTrain(w http.ResponseWriter, r *http.Request) {
	_, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}

	query := r.URL.Query()
	dateStr := query.Get("date")
	source := query.Get("source")
	destination := query.Get("destination")

	var filterDate time.Time

	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			logger.LogToFile("Invalid date format")
			response.ErrorResponse(w, http.StatusBadRequest, "Invalid date", 1001)
			return
		}
		filterDate = parsedDate
	}

	trains, err := h.TrainService.GetTrains(filterDate, source, destination)

	if err != nil {
		logger.LogToFile("error in fetching train details")
		response.ErrorResponse(w, http.StatusUnauthorized, "error in fetching train details", 1007)
		return
	}

	response.SuccessResponse(w, trains, "Train details fetched successfully", http.StatusOK)
}
