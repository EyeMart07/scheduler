package api

import "github.com/gin-gonic/gin"

func (a *App) RegisterEndpoints(router *gin.Engine) {
	router.GET("/users", a.GetUsers)
	router.GET("/users/:id", a.GetUsersById)

	router.POST("/appointments", a.CreateAppointment)
	router.GET("/appointments/:id", a.GetAppointmentById)
	router.GET("/appointments", a.GetAppointments)

	router.POST("/cars", a.AddCar)
	router.GET("/cars", a.GetCars)
	router.GET("/cars/:customer", a.GetCarsByCustomer)
	router.GET("/cars/:customer/:id", a.GetCarByCustomerAndId)
}
