package app

import "github.com/gin-gonic/gin"

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

// StartApp ... maps
func StartApp() {
	mapUrls()

	if err := router.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
