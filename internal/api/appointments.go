package api

import (
	"net/http"

	"github.com/EyeMart07/scheduler/internal/store"
	"github.com/gin-gonic/gin"
)

type AppointmentReqs struct {
	DateTime string `json:"appointment_date_time"`
	Customer int    `json:"customer_id"`
	Notes    string `json:"notes"`
}

func (a *App) GetAppointmentById(c *gin.Context) {
	id := c.Param("id")

	app, err := a.Store.GetAppointmentById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
	}

	c.IndentedJSON(http.StatusOK, app)

}

func (a *App) GetAppointments(c *gin.Context) {
	app, err := a.Store.GetAppointments()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
	}

	c.IndentedJSON(http.StatusOK, app)

}

// adds the appointment associates it a user, and adds it to the database
func (a *App) CreateAppointment(c *gin.Context) {
	var newApp AppointmentReqs

	// gets the appointment data from the request
	if err := c.BindJSON(&newApp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing appointment data"})
		return
	}

	if err := a.Store.CreateAppointment(store.Appointment{
		DateTime: newApp.DateTime,
		Customer: newApp.Customer,
		Notes:    newApp.Notes,
	}); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		return
	}

	// return the created status
	c.IndentedJSON(http.StatusCreated, newApp)
}
