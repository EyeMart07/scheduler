package api

import (
	"net/http"

	"github.com/EyeMart07/scheduler/internal/store"
	"github.com/gin-gonic/gin"
)

type Availability struct {
	Date  string `json:"date"`
	Start string `json:"start_time"`
	End   string `json:"end_time"`
}

func (a *App) GetAvailability(c *gin.Context) {

	date := c.Query("date")

	avail, err := a.Store.GetAvailability(store.AvailabilityArguments{
		Date: &date,
	})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, avail)
}

func (a *App) ChangeAvailability(c *gin.Context) {
	date := c.Param("date")

	var changes Availability

	if err := c.BindJSON(&changes); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid availability data"})
		return
	}

	if err := a.Store.ChangeAvailability(store.Availability{
		Date:  date,
		Start: changes.Start,
		End:   changes.End,
	}); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully updated availability"})
}

func (a *App) GetTimeSlots(c *gin.Context) {
	date := c.Query("date")

	timeSlots, err := a.Store.GetTimeSlots(store.TimeSlotArguments{
		Date: &date,
	})

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, timeSlots)

}

func (a *App) AddAvailability(c *gin.Context) {
	var newAvailability Availability

	if err := c.BindJSON(&newAvailability); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid availability data"})
		return
	}

	if err := a.Store.AddAvailability(store.Availability{
		Date:  newAvailability.Date,
		Start: newAvailability.Start,
		End:   newAvailability.End,
	}); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating availability"})
		return
	}

	// return the created status
	c.IndentedJSON(http.StatusCreated, newAvailability)

}
