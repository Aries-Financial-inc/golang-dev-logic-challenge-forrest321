package routes

import (
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-forrest321/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/analyze", controllers.AnalysisHandler)
	return router
}
