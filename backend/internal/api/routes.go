package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) RegisterEndpoints(router *gin.Engine) {

	router.POST("/appointments", a.CreateAppointment)
	router.DELETE("/appointments/:id", a.DeleteAppointment)
	router.PUT("/appointments/:id", a.ChangeAppointment)
	router.GET("/admin/appointments", a.CheckAuth, a.CheckAdmin, a.GetAppointments)

	router.PUT("/admin/availability/:date", a.CheckAuth, a.CheckAdmin, a.ChangeAvailability)
	router.POST("/admin/availability", a.CheckAuth, a.CheckAdmin, a.AddAvailability)
	router.GET("/availability", a.CheckAuth, a.GetAvailability)
	router.GET("/timeslots", a.CheckAuth, a.GetTimeSlots)

	router.POST("/register", a.SignUp)
	router.POST("/signin", a.SignIn)

	router.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "ok"})
	})
}
