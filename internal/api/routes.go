package api

import "github.com/gin-gonic/gin"

func (a *App) RegisterEndpoints(router *gin.Engine) {

	router.POST("/appointments", a.CreateAppointment)
	router.GET("/appointments/:id", a.GetAppointmentById)
	router.GET("/appointments", a.GetAppointments)
}
