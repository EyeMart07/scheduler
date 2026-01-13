package api

import "github.com/gin-gonic/gin"

func (a *App) RegisterEndpoints(router *gin.Engine) {

	router.POST("/admin/appointments", a.CheckAuth, a.CheckAdmin, a.CreateAppointment)
	router.GET("/admin/appointments", a.CheckAuth, a.CheckAdmin, a.GetAppointments)

	router.POST("/admin/availability", a.CheckAuth, a.CheckAdmin, a.AddAvailability)
	router.GET("/availability", a.CheckAuth, a.GetAvailability)
	router.GET("/timeslots", a.CheckAuth, a.GetTimeSlots)

	router.POST("/register", a.SignUp)
	router.POST("/signin", a.SignIn)
}
