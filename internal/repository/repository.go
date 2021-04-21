package repository

import (
	"time"

	"github.com/winadiw/go-bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestrictions(r models.RoomRestriction) error

	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomID int) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
}
