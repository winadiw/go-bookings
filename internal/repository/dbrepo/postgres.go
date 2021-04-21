package dbrepo

import (
	"context"
	"time"

	"github.com/winadiw/go-bookings/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestrictions inserts a room restrictions into the database
func (m *postgresDBRepo) InsertRoomRestrictions(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		time.Now(),
		time.Now(),
		res.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDates returns true if availability exists for roomID, false if no availability
func (m *postgresDBRepo) SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `SELECT count(id)
	FROM room_restrictions
	WHERE room_id = $1 AND $2 < end_date AND $3 > start_date `

	err := m.DB.QueryRowContext(ctx, query, roomID, start, end).Scan(&numRows)

	if err != nil {
		return false, err
	}

	// Must have no data
	return numRows == 0, nil
}
