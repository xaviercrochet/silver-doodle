package controllers

import (
	"localsearch-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPlace ...
func GetPlace(c *gin.Context) {
	placeID := c.Param("place_id")
	place, err := services.PlacesService.GetPlace(placeID)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, place)
}
