package api

import "github.com/gin-gonic/gin"

func (a *App) RegisterEndpoints(router *gin.Engine) {

	router.POST("/appointments", a.CreateAppointment)
	router.GET("/appointments/:date", a.GetAppointmentOnDay)
	router.GET("/appointments", a.GetAppointments)

	router.POST("/availability", a.AddAvailability)
	router.GET("/availability/:date", a.GetAvailability)
	router.GET("/timeslots/:date", a.GetTimeSlots)
}
