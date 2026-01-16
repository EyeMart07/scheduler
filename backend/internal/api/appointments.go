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

/*
Returns appointments on a given day

accepted queries:
fromDate, toDate : enables queries on a range of dates
*/
func (a *App) GetAppointments(c *gin.Context) {
	toDate := c.Query("to_date")
	fromDate := c.Query("to")
	app, err := a.Store.GetAppointments(store.AppointmentArguments{
		FromDate: &fromDate,
		ToDate:   &toDate,
	})

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, app)

}

func (a *App) DeleteAppointment(c *gin.Context) {
	id := c.Param("id")

	email, phone, err := a.Store.DeleteAppointment(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error deleting appointment"})
		return
	}
	if email == "" && phone == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "appointment doesn't exist"})
		return
	}

	// SEND CONFIRMATION HERE
	c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

func (a *App) ChangeAppointment(c *gin.Context) {
	id := c.Param("id")
	var changes AppointmentReqs
	// gets the appointment data from the request
	if err := c.BindJSON(&changes); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid appointment data"})
		return
	}

	email, phone, err := a.Store.ChangeAppointment(id, store.Appointment{
		Notes:     changes.Notes,
		FirstName: changes.FirstName,
		LastName:  changes.LastName,
		Email:     changes.Email,
		Phone:     changes.Phone,
		Address:   changes.Address,
		Make:      changes.Make,
		Model:     changes.Model,
		Year:      changes.Year,
		Vin:       changes.Vin,
		Mileage:   changes.Mileage,
		Date:      changes.Date,
		Start:     changes.Start,
		End:       changes.End,
	})

	if err != nil || (email == "" && phone == "") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error updating appointment"})
		return
	}

	// SEND CONFIRMATION HERE
	c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully updated"})
}

func (a *App) CreateAppointment(c *gin.Context) {
	var newApp AppointmentReqs

	// gets the appointment data from the request
	if err := c.BindJSON(&newApp); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid appointment data"})
		return
	}

	id, err := a.Store.CreateAppointment(store.Appointment{
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
	})

	if id == "" || err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		return
	}
	// return the created status
	c.IndentedJSON(http.StatusCreated, newApp)
}
