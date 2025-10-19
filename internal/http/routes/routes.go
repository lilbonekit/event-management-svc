package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/internal/handlers"
)

func SetupUserRoutes(r *gin.Engine, uh *handlers.UserHandler) {
	r.POST("/signup", uh.SignUp)
	r.POST("/login", uh.SignIn)
}

func SetupEventRoutes(r *gin.Engine, eh *handlers.EventHandler, auth gin.HandlerFunc) {
	r.GET("/events", eh.List)
	r.GET("/events/:id", eh.Get)

	// protected routes
	r.POST("/events", auth, eh.Create)
	r.PUT("/events/:id", auth, eh.Update)
	r.DELETE("/events/:id", auth, eh.Delete)

	r.POST("/events/:id/registrations", auth, eh.Register)
	r.DELETE("/events/:id/registrations", auth, eh.CancelRegistration)
}
