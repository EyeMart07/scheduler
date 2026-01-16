package store

import (
	"fmt"
	"strings"
	"sync"
)

type Availability struct {
	Date  string `json:"date"`
	Start string `json:"start_time"`
	End   string `json:"end_time"`
}

type AvailabilityArguments struct {
	Date *string `json:"date"`
}

func buildAvailabilityQuery(queried AvailabilityArguments) (string, []any) {
	q := `SELECT date, start_time, end_time FROM availability`

	where := []string{}
	args := []any{}

	add := func(cond string, val any) {
		args = append(args, val)
		where = append(where, fmt.Sprintf(cond, len(args))) // %d becomes $1, $2...
	}

	// date = matches a whole day (recommended: [dayStart, nextDay))
	if queried.Date != nil {
		add("date = $%d", *queried.Date)
	}

	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}

	return q, args
}

func (s *Store) GetAvailability(queried AvailabilityArguments) (Availability, error) {
	query, args := buildAvailabilityQuery(queried)
	row := s.DB.QueryRow(query, args...)

	var avail Availability

	if err := row.Scan(&avail.Date, &avail.Start, &avail.End); err != nil {
		return avail, err
	}

	return avail, nil
}

type TimeSlot struct {
	Start string `json:"start_time"`
	End   string `json:"end_time"`
}

type TimeSlotArguments struct {
	Date *string `json:"date"`
}

func (s *Store) ChangeAvailability(changes Availability) error {
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	// formats the query with the given data
	_, err = tx.Exec("UPDATE availability SET start_time=$1, end_time=$2 WHERE date=$3", changes.Start, changes.End, changes.Date)

	if err != nil {
		tx.Rollback()
		return err
	}

	// attempts to commit the transaction if the query succeeds
	err = tx.Commit()

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetTimeSlots(queried TimeSlotArguments) ([]TimeSlot, error) {
	var availability Availability
	var appointments []Appointment
	var err1, err2 error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		availability, err1 = s.GetAvailability(AvailabilityArguments{
			Date: queried.Date,
		})
	}()

	go func() {
		defer wg.Done()
		appointments, err2 = s.GetAppointments(AppointmentArguments{
			ToDate: queried.Date,
		})
	}()

	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	//if no appointments for that return one time slot
	if len(appointments) == 0 {
		return []TimeSlot{{Start: availability.Start, End: availability.End}}, nil
	}
	// parse through all the apointments for that day creating all the available time slots
	var availTimes []TimeSlot
	curStart := availability.Start
	for _, app := range appointments {
		availTimes = append(availTimes, TimeSlot{
			Start: curStart,
			End:   app.Start,
		})
		curStart = app.End
	}

	availTimes = append(availTimes, TimeSlot{
		Start: curStart,
		End:   availability.End,
	})

	return availTimes, nil
}

func (s *Store) AddAvailability(avail Availability) error {
	// begins a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	// attempts to execute the query
	_, err = tx.Exec("INSERT INTO availability(date, start_time, end_time) VALUES($1, $2, $3)", avail.Date, avail.Start, avail.End)

	if err != nil {
		tx.Rollback()
		return err
	}

	// attempts to commit the transaction if the query succeeds
	err = tx.Commit()

	if err != nil {
		return err
	}
	return nil
}
