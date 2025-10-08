package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/middlewares"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/events", getEvents)
	router.GET("/events/:id", getEvent)

	authenticated := router.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.POST("/events/:id/registrations", registerForEvent)

	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.DELETE("/events/:id/registrations", cancelRegistration)

	router.POST("/signup", signUp)
	router.POST("/login", signIn)
}
