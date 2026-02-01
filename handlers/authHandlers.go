package handlers

import (
	"net/http"
	"strconv"
	"time"
	"workout-tracker/config"
	"workout-tracker/models"
	"workout-tracker/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := handler.authRepo.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (handler *AuthHandler) SignIn(c *gin.Context) {
	var request models.SignInRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid parameters"))
		return
	}
	user, err := handler.authRepo.FindByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Can not find by email"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid credials"))
		return
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not sign JWT"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
