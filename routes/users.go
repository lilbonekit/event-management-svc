package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svg/models"
	"github.com/lilbonekit/event-management-svg/utils"
)

// signUp godoc
// @Summary      Create a new user (signup)
// @Description  Registers a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body     models.User true "User signup"
// @Success      201  {object} map[string]interface{}
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /signup [post]
func signUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if err := user.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
		return
	}

	// Never return password in API responses
	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user_id": user.ID})
}

// signIn godoc
// @Summary      Sign in and receive a JWT
// @Description  Authenticates user with email and password; returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.User true "User credentials"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /login [post]
func signIn(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if err := user.ValidateCredentials(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User signed in", "token": token})
}
