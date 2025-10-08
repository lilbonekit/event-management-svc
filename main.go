package main

// @title           Event Management REST API
// @version         1.0
// @description     Simple event management API built with Gin and JWT authentication.
// @BasePath        /

/*
@securityDefinitions.apikey BearerAuth
@in header
@name Authorization
@description Use format: Bearer <JWT token>
*/
import (
	"github.com/joho/godotenv"

	_ "github.com/lilbonekit/event-management-svc/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/db"
	"github.com/lilbonekit/event-management-svc/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(".env file not found")
	}

	db.InitDB()
	server := gin.Default()

	routes.SetupRoutes(server)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Run(":8089")
}
