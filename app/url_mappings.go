package app

import "localsearch-api/controllers"

func mapUrls() {
	router.GET("/places/:place_id", controllers.GetPlace)
}
