package store

import "fmt"

// car structure
type Car struct {
	Make     string `json:"make"`
	Model    string `json:"model"`
	Year     int    `json:"year"`
	Mileage  int    `json:"mileage"`
	Vin      string `json:"vin"`
	Customer int    `json:"customer"`
}

// returns a specific car connected to a customer
func (s *Store) GetCarByCustomerAndId(customer string, id string) (Car, error) {
	row := s.DB.QueryRow("SELECT (make, model, year, mileage, vin, customer) FROM cars WHERE customer=$1 and id=$2", customer, id)
	var car Car
	err := row.Scan(&car.Make, &car.Model, &car.Year, &car.Mileage, &car.Vin, &car.Customer)
	return car, err
}

// returns all cars connected to a user
func (s *Store) GetCarsByCustomer(customer string) ([]Car, error) {
	rows, err := s.DB.Query("SELECT (make, model, year, mileage, vin, customer)) FROM cars where customer=$1", customer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cars := []Car{}
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.Make, &car.Model, &car.Year, &car.Mileage, &car.Vin, &car.Customer); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

// returns all cars (mainly for testing purposes)
func (s *Store) GetCars() ([]Car, error) {
	rows, err := s.DB.Query("SELECT (make, model, year, mileage, vin, customer) FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cars := []Car{}
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.Make, &car.Model, &car.Year, &car.Mileage, &car.Vin, &car.Customer); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func (s *Store) AddCar(car Car) error {
	// begins a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	// formats the query with the given data
	query := fmt.Sprintf("INSERT INTO cars(model, make, year, mileage, vin, customer) VALUES('%s', '%s', %d, %d, '%s', %d)", car.Make, car.Model, car.Year, car.Mileage, car.Vin, car.Customer)

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
