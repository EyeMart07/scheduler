package store

import (
	"fmt"
	"sync"
)

type Availability struct {
	Date  string `json:"date"`
	Start string `json:"start_time"`
	End   string `json:"end_time"`
}

func (s *Store) GetAvailability(day string) (Availability, error) {
	row := s.DB.QueryRow("SELECT date, start_time, end_time FROM availability WHERE date=$1", day)

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

func (s *Store) GetTimeSlots(day string) ([]TimeSlot, error) {
	var availability Availability
	var appointments []Appointment
	var err1, err2 error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		availability, err1 = s.GetAvailability(day)
	}()

	go func() {
		defer wg.Done()
		appointments, err2 = s.GetAppointmentOnDay(day)
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

	// formats the query with the given data
	query := fmt.Sprintf("INSERT INTO availability(date, start_time, end_time) VALUES('%s', '%s', '%s')", avail.Date, avail.Start, avail.End)
	// attempts to execute the query
	_, err = tx.Exec(query)

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
