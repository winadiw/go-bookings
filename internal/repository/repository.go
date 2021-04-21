package repository

import "github.com/winadiw/go-bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestrictions(r models.RoomRestriction) error
}
