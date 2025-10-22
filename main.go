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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lilbonekit/event-management-svc/internal/db"
	"github.com/lilbonekit/event-management-svc/internal/handlers"
	"github.com/lilbonekit/event-management-svc/internal/http/routes"
	"github.com/lilbonekit/event-management-svc/internal/middlewares"
	"github.com/lilbonekit/event-management-svc/internal/repo"
	"github.com/lilbonekit/event-management-svc/internal/service"
)

func main() {
	_ = godotenv.Load()

	db, err := db.Open(db.Config{
		Path:         "api.db",
		PingTimeout:  2 * time.Second,
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repo.NewUserRepo(db)
	userSvc := service.NewUserService(db, userRepo)
	userH := handlers.NewUserHandler(userSvc)

	r := gin.Default()
	routes.SetupUserRoutes(r, userH)

	evRepo := repo.NewEventRepo(db)
	evSvc := service.NewEventService(db, evRepo)
	evH := handlers.NewEventHandler(evSvc)

	routes.SetupEventRoutes(r, evH, middlewares.Authenticate)

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("No PORT environment variable detected, defaulting to 8080")
		port = "8080"
	}

	r.Run(":" + port)
}
