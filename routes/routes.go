package routes

import (
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-forrest321/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRouter builds the routes for the service. Since this should only set up routes,
// validation logic was moved to the handler
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/analyze", controllers.AnalysisHandler)
	return router
}
