package store

import (
	"database/sql"
	"fmt"
	"strings"
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

type AppointmentArguments struct {
	Date *string `json:"date"`
}

func buildAppointmentQuery(queried AppointmentArguments) (string, []any) {
	q := `SELECT notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage, date, start_time, end_time FROM appointments`

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

	q += " ORDER BY start_time ASC"

	return q, args
}

func (s *Store) GetAppointments(queried AppointmentArguments) ([]Appointment, error) {

	query, args := buildAppointmentQuery(queried)

	rows, err := s.DB.Query(query, args...)

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
	_, err = tx.Exec("INSERT INTO appointments(notes, first_name, last_name, email, phone, address, make, model, year, vin, mileage, date, start_time, end_time) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", app.Notes, app.FirstName, app.LastName, app.Email, app.Phone, app.Address, app.Make, app.Model, app.Year, app.Vin, app.Mileage, app.Date, app.Start, app.End)

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
