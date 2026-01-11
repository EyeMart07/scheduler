package store

import (
	"database/sql"
	"fmt"
)

type Appointment struct {
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
	Date      string `json:"date"`
	Start     string `json:"start_time"`
	End       string `json:"end_time"`
}

func parseRows(rows *sql.Rows) ([]Appointment, error) {
	defer rows.Close()

	appointments := []Appointment{}
	for rows.Next() {
		var app Appointment
		if err := rows.Scan(&app.Notes, &app.FirstName, &app.LastName, &app.Email, &app.Phone, &app.Address, &app.Make, &app.Model, &app.Year, &app.Vin, &app.Mileage, &app.Date, &app.Start, &app.End); err != nil {
			return appointments, err
		}
		appointments = append(appointments, app)
	}
	return appointments, nil
}

func (s *Store) GetAppointmentOnDay(date string) ([]Appointment, error) {
	rows, err := s.DB.Query("SELECT notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage, date, start_time, end_time FROM appointments WHERE date=$1 ORDER BY start_time", date)

	if err != nil {
		return nil, err
	}

	return parseRows(rows)
}

func (s *Store) GetAppointments() ([]Appointment, error) {
	rows, err := s.DB.Query("SELECT notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage, date, start_time, end_time FROM appointments")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseRows(rows)
}

func (s *Store) CreateAppointment(app Appointment) error {
	// begins a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	// formats the query with the given data
	query := fmt.Sprintf("INSERT INTO appointments(notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage, date, start_time, end_time) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, '%s', %d, '%s', '%s', '%s')", app.Notes, app.FirstName, app.LastName, app.Email, app.Phone, app.Address, app.Make, app.Model, app.Year, app.Vin, app.Mileage, app.Date, app.Start, app.End)
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
