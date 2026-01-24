package handlers

import (
	"net/http"
	"workout-tracker/models"
	"workout-tracker/repositories"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authRepo *repositories.AuthRepository
}

func NewAuthHandler(authRepo *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{authRepo: authRepo}
}

func (handler *AuthHandler) SignUp(c *gin.Context) {
	var request models.RegisterUserRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not bind json body"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not hash password"))
		return
	}
}
