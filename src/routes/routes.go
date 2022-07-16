package routes

import (
	"hexdecimal16/chaipay-assignment/src/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/api/v1/health", controllers.Health)
	router.GET("/api/v1/get_intents", controllers.GetPaymentIntents)
	router.POST("/api/v1/create_intent", controllers.CreateIntent)
	router.POST("/api/v1/capture_intent/:id", controllers.CaptureIntent)
	router.POST("/api/v1/create_refund/:id", controllers.CreateRefund)
	router.NoRoute(controllers.Error404)
}
