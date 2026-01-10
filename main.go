package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// server structure that will be used in handlers
type Server struct {
	DB *sql.DB
}

func main() {
	// loads the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// gets the connection string from the .env file and creates a connection to the database
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	// closes the connection once the program ends
	defer db.Close()

	// ensures the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// s is a pointer to a server struct so the same connection can be passed to each endpoint handler
	s := &Server{DB: db}

	// sets up the server
	router := gin.Default()

	router.GET("/users", s.getUsers) // we pass the get users function that is referenced with s
	router.GET("/users/:id", s.getUsersById)
	router.POST("/appointments", s.createAppointment)
	router.GET("/appointments/:id", s.getAppointment)
	router.GET("/appointments", s.getAppointments)

	router.Run()
}

// user structure
type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

type Appointment struct {
	DateTime string `json:"appointment_date_time"`
	Customer int    `json:"customer_id"`
	Notes    string `json:"notes"`
}

func (s *Server) getAppointment(c *gin.Context) {
	id := c.Param("id")
	row := s.DB.QueryRow("SELECT appointment_date_time, customer_id, notes FROM appointments WHERE customer_id=$1", id)

	var app Appointment

	if err := row.Scan(&app.DateTime, &app.Customer, &app.Notes); err != nil {
		c.IndentedJSON(http.StatusNotFound, app)
		return
	}

	c.IndentedJSON(http.StatusOK, app)

}

func (s *Server) getAppointments(c *gin.Context) {
	rows, err := s.DB.Query("SELECT appointment_date_time, customer_id, notes FROM appointments")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	appointments := []Appointment{}

	for rows.Next() {
		var app Appointment

		if err := rows.Scan(&app.DateTime, &app.Customer, &app.Notes); err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "scan failed"})
			return
		}

		appointments = append(appointments, app)
	}

	c.IndentedJSON(http.StatusOK, appointments)

}

// adds the appointment associates it a user, and adds it to the database
func (s *Server) createAppointment(c *gin.Context) {
	var newApp Appointment

	if err := c.BindJSON(&newApp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing appointment data"})
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		tx.Rollback()
		return
	}

	query := fmt.Sprintf("INSERT INTO appointments(appointment_date_time, customer_id, notes) VALUES('%s'::timestamptz, %d, '%s')", newApp.DateTime, newApp.Customer, newApp.Notes)

	_, err = tx.Exec(query)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		tx.Rollback()
		return
	}

	c.IndentedJSON(http.StatusCreated, newApp)
}

// returns a specific user by id
func (s *Server) getUsersById(c *gin.Context) {
	id := c.Param("id")
	row := s.DB.QueryRow("SELECT * FROM users WHERE id=$1", id)

	var user User

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Address); err != nil {
		c.IndentedJSON(http.StatusNotFound, user)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

// returns all users (mainly for testing purposes)
func (s *Server) getUsers(c *gin.Context) {
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
