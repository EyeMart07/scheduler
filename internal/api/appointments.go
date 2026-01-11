package api

import (
	"net/http"

	"github.com/EyeMart07/scheduler/internal/store"
	"github.com/gin-gonic/gin"
)

type AppointmentReqs struct {
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

func (a *App) GetAppointmentOnDay(c *gin.Context) {
	date := c.Param("date")

	app, err := a.Store.GetAppointmentOnDay(date)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, app)

}

func (a *App) GetAppointments(c *gin.Context) {
	app, err := a.Store.GetAppointments()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, app)

}

func (a *App) CreateAppointment(c *gin.Context) {
	var newApp AppointmentReqs

	// gets the appointment data from the request
	if err := c.BindJSON(&newApp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid appointment data"})
		return
	}

	if err := a.Store.CreateAppointment(store.Appointment{
		Notes:     newApp.Notes,
		FirstName: newApp.FirstName,
		LastName:  newApp.LastName,
		Email:     newApp.Email,
		Phone:     newApp.Phone,
		Address:   newApp.Address,
		Make:      newApp.Make,
		Model:     newApp.Model,
		Year:      newApp.Year,
		Vin:       newApp.Vin,
		Mileage:   newApp.Mileage,
		Date:      newApp.Date,
		Start:     newApp.Start,
		End:       newApp.End,
	}); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		return
	}

	// return the created status
	c.IndentedJSON(http.StatusCreated, newApp)
}
