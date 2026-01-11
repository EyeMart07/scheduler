package api

import (
	"net/http"

	"github.com/EyeMart07/scheduler/internal/store"
	"github.com/gin-gonic/gin"
)

type CarReq struct {
	Make     string `json:"make"`
	Model    string `json:"model"`
	Year     int    `json:"year"`
	Mileage  int    `json:"mileage"`
	Vin      string `json:"vin"`
	Customer int    `json:"customer"`
}

func (a *App) GetCarByCustomerAndId(c *gin.Context) {
	id := c.Param("id")
	customer := c.Param("customer")

	car, err := a.Store.GetCarByCustomerAndId(customer, id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, car)

}

func (a *App) GetCarsByCustomer(c *gin.Context) {
	customer := c.Param("customer")
	cars, err := a.Store.GetCarsByCustomer(customer)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, cars)

}

func (a *App) GetCars(c *gin.Context) {
	cars, err := a.Store.GetCars()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, cars)

}

// adds the appointment associates it a user, and adds it to the database
func (a *App) AddCar(c *gin.Context) {
	var newCar CarReq

	// gets the appointment data from the request
	if err := c.BindJSON(&newCar); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing appointment data"})
		return
	}

	if err := a.Store.AddCar(store.Car{
		Make:     newCar.Make,
		Model:    newCar.Model,
		Year:     newCar.Year,
		Mileage:  newCar.Mileage,
		Vin:      newCar.Vin,
		Customer: newCar.Customer,
	}); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "error creating appointment"})
		return
	}

	// return the created status
	c.IndentedJSON(http.StatusCreated, newCar)
}
