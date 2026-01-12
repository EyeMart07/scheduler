package store

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password_hash"`
}

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password_hash"`
}

type AuthResp struct {
	Customer string `json:"customer"`
	Password string `json:"password_hash"`
}

// returns the user id if the log in is successful
func (s *Store) Authorize(credentials AuthReq) string {

	var user AuthResp
	row := s.DB.QueryRow("SELECT customer password_hash FROM users WHERE email=$1", credentials.Email)

	if err := row.Scan(&user); err != nil {
		return ""
	}

	attempt, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)
	if err := bcrypt.CompareHashAndPassword(attempt, []byte(user.Password)); err != nil {
		return ""
	}

	return user.Customer
}

// cretes a user and returns the associated id in the table
func (s *Store) CreateUser(user User) (string, error) {

	// begins a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// formats the query with the given data
	query := fmt.Sprintf("INSERT INTO users(first_name, last_name, email, password_hash) VALUES('%s', '%s', '%s', '%s')", user.FirstName, user.LastName, user.Email, user.Password)
	// attempts to execute the query
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	row := tx.QueryRow("SELECT id from users where email=%s", user.Email)

	var customerId string
	if err := row.Scan(&customerId); err != nil {
		tx.Rollback()
		return "", err
	}

	// attempts to commit the transaction if the query succeeds
	err = tx.Commit()

	if err != nil {
		return "", err
	}
	return customerId, nil
}
