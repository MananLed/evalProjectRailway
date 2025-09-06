package model

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID          uuid.UUID `json:"id"`
	TrainID     uuid.UUID `json:"trainid"`
	PassengerID uuid.UUID `json:"passengerid"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Departure   time.Time `json:"departure"`
	TotalFare   int       `json:"totalfare"`
	BookedSeats int       `json:"bookedseats"`
}
