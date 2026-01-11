package store

import (
	"fmt"
)

type Appointment struct {
	DateTime  string `json:"appointment_date_time"`
	Notes     string `json:"notes"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Make      string `json:"make"`
	Model     string `json:"model"`
	Year      int    `json:"year"`
	Vin       string `json:"vin"`
	Mileage   int    `json:"mileage"`
}

func (s *Store) GetAppointmentById(id string) (Appointment, error) {
	row := s.DB.QueryRow("SELECT appointment_date_time, notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage FROM appointments WHERE customer_id=$1", id)

	var app Appointment

	if err := row.Scan(&app.DateTime, &app.Notes, &app.FirstName, &app.LastName, &app.Email, &app.Phone, &app.Address, &app.Make, &app.Model, &app.Year, &app.Vin, &app.Mileage); err != nil {
		return app, err
	}

	return app, nil
}

func (s *Store) GetAppointments() ([]Appointment, error) {
	rows, err := s.DB.Query("SELECT appointment_date_time, notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage FROM appointments")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := []Appointment{}
	for rows.Next() {
		var app Appointment
		if err := rows.Scan(&app.DateTime, &app.Notes, &app.FirstName, &app.LastName, &app.Email, &app.Phone, &app.Address, &app.Make, &app.Model, &app.Year, &app.Vin, &app.Mileage); err != nil {
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
	query := fmt.Sprintf("INSERT INTO appointments(appointment_date_time, notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage) VALUES('%s'::timestamptz, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, '%s', %d)", app.DateTime, app.Notes, app.FirstName, app.LastName, app.Email, app.Phone, app.Address, app.Make, app.Model, app.Year, app.Vin, app.Mileage)

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
