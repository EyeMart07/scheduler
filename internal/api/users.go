package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// returns a specific user by id
func (a *App) getUsersById(c *gin.Context) {
	id := c.Param("id")

	user, err := a.Store.getUsersById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
	}

	c.IndentedJSON(http.StatusOK, user)
}
