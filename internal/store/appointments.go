package store

import (
	"fmt"
)

type Appointment struct {
	DateTime string `json:"appointment_date_time"`
	Customer int    `json:"customer_id"`
	Notes    string `json:"notes"`
}

func (s *Store) GetAppointmentById(id string) (Appointment, error) {
	row := s.DB.QueryRow("SELECT appointment_date_time, customer_id, notes FROM appointments WHERE customer_id=$1", id)

	var app Appointment

	if err := row.Scan(&app.DateTime, &app.Customer, &app.Notes); err != nil {
		return app, err
	}

	return app, nil
}

func (s *Store) GetAppointments() ([]Appointment, error) {
	rows, err := s.DB.Query("SELECT appointment_date_time, customer_id, notes FROM appointments")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := []Appointment{}
	for rows.Next() {
		var app Appointment
		if err := rows.Scan(&app.DateTime, &app.Customer, &app.Notes); err != nil {
			return appointments, err
		}
		appointments = append(appointments, app)
	}
	return appointments, err
}

// adds the appointment associates it a user, and adds it to the database
func (s *Store) CreateAppointment(app Appointment) error {
	// begins a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	// formats the query with the given data
	query := fmt.Sprintf("INSERT INTO appointments(appointment_date_time, customer_id, notes) VALUES('%s'::timestamptz, %d, '%s')", app.DateTime, app.Customer, app.Notes)

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
