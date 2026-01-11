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

	day := c.Param("date")

	avail, err := a.Store.GetAvailability(day)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, avail)
}

func (a *App) GetTimeSlots(c *gin.Context) {
	day := c.Param("date")

	timeSlots, err := a.Store.GetTimeSlots(day)

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
