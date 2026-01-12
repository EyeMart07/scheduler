package api

import "github.com/gin-gonic/gin"

func (a *App) RegisterEndpoints(router *gin.Engine) {

	router.POST("/appointments", a.CheckAuth, a.CreateAppointment)
	router.GET("/admin/appointments/:date", a.CheckAuth, a.GetAppointmentOnDay)
	router.GET("/admin/appointments", a.CheckAuth, a.GetAppointments)

	router.POST("/admin/availability", a.CheckAuth, a.AddAvailability)
	router.GET("/availability/:date", a.CheckAuth, a.GetAvailability)
	router.GET("/timeslots/:date", a.CheckAuth, a.GetTimeSlots)

	router.POST("/register", a.SignUp)
	router.POST("/signin", a.SignIn)
}
