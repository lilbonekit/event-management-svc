package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/internal/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid token"})
		return
	}

	userID, err := utils.ValidateToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	context.Set("userID", int64(userID))
	context.Next()
}
