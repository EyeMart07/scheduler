package store

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// user structure
type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

// returns a specific user by id
func (s *Store) getUsersById(id int) (User, error) {
	row := s.DB.QueryRow("SELECT * FROM users WHERE id=$1", id)

	var user User

	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Address)

	return user, err
}

// returns all users (mainly for testing purposes)
func (s *Store) getUsers(c *gin.Context) {
	rows, err := s.DB.Query("SELECT * FROM users")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Address); err != nil {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "scan failed"})
			return
		}

		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}
