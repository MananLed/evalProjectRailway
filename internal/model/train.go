package model

import (
	"time"

	"github.com/google/uuid"
)

type Train struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Departure   time.Time `json:"departure"`
	Arrival     time.Time `json:"arrival"`
	TotalSeats  int       `json:"totalseats"`
	BookedSeats int       `json:"bookedseats"`
	SeatFare    int       `json:"seatfare"`
}
