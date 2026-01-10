package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// returns a specific user by id
func (a *App) GetUsersById(c *gin.Context) {
	id := c.Param("id")
	user, err := a.Store.GetUsersById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
	c.IndentedJSON(http.StatusOK, user)
}

// returns a specific user by id
func (a *App) GetUsers(c *gin.Context) {
	users, err := a.Store.GetUsers()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
	c.IndentedJSON(http.StatusOK, users)
}
