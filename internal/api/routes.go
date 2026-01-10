package api

import "github.com/gin-gonic/gin"

func (a *App) RegisterEndpoints(router *gin.Engine) {
	router.GET("/users", a.GetUsers) // we pass the get users function that is referenced with s
	router.GET("/users/:id", a.GetUsersById)
	router.POST("/appointments", a.CreateAppointment)
	router.GET("/appointments/:id", a.GetAppointmentById)
	router.GET("/appointments", a.GetAppointments)
}
