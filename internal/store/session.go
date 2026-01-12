package store

import (
	"crypto/sha256"
	"fmt"
)

type Session struct {
	Customer    string   `json:"customer"`
	SessionHash [32]byte `json:"session_hash"`
}

func (s *Store) CheckAuth(token string) string {
	hashed := sha256.Sum256([]byte(token))

	row := s.DB.QueryRow("SELECT customer FROM sessions WHERE session_hash=$1", hashed)
	var customer string
	if err := row.Scan(&customer); err != nil {
		return ""
	}
	return customer
}

func (s *Store) CreateSession(session Session) error {

	// begins a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	// formats the query with the given data
	query := fmt.Sprintf("INSERT INTO appointments(customer, session_hash) VALUES('%s', '%b')", session.Customer, session.SessionHash)
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
