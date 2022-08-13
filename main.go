package main

import (
	"log"

	gin "github.com/gin-gonic/gin"

	"github.com/hexdecimal16/chaipay-assignment/database"
	"github.com/hexdecimal16/chaipay-assignment/src/models"
	"github.com/hexdecimal16/chaipay-assignment/src/routes"
)

func main() {
	// Database
	database.ConnectDB()

	// Migrate the schema
	database.DB.AutoMigrate(&models.PaymentIntent{})

	// Init Router
	router := gin.Default()

	// Route Handlers / Endpoints
	routes.Routes(router)

	log.Fatal(router.Run(":5000"))
}
