package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EyeMart07/scheduler/internal/api"
	"github.com/EyeMart07/scheduler/internal/db"
	"github.com/EyeMart07/scheduler/internal/store"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	db, err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}
	// closes the connection once the program ends
	defer db.Close()

	// s is a pointer to a server struct so the same connection can be passed to each endpoint handler
	st := store.NewDatabase(db)

	app := api.NewApp(st)

	// sets up the server
	router := gin.Default()

	router.GET("/users", s.getUsers) // we pass the get users function that is referenced with s
	router.GET("/users/:id", s.getUsersById)
	router.POST("/appointments", s.createAppointment)
	router.GET("/appointments/:id", s.getAppointment)
	router.GET("/appointments", s.getAppointments)

	router.Run()
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

	// gets the appointment data from the request
	if err := c.BindJSON(&newApp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing appointment data"})
		return
	}

	// begins a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		tx.Rollback()
		return
	}

	// formats the query with the given data
	query := fmt.Sprintf("INSERT INTO appointments(appointment_date_time, customer_id, notes) VALUES('%s'::timestamptz, %d, '%s')", newApp.DateTime, newApp.Customer, newApp.Notes)

	// attempts to execute the query
	_, err = tx.Exec(query)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		tx.Rollback()
		return
	}

	// attempts to commit the transaction if the query succeeds
	err = tx.Commit()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		tx.Rollback()
		return
	}

	// return the created status
	c.IndentedJSON(http.StatusCreated, newApp)
}
