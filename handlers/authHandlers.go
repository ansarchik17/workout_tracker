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

	err := c.ShouldBind(&request)
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

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("invalid request body"))
		return
	}

	user, err := handler.authRepo.FindByEmail(
		c.Request.Context(),
		request.Email,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("invalid credentials"))
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("invalid credentials"))
		return
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.UserID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
	}

	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	).SignedString([]byte(config.Config.JwtSecretKey))

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not sign token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
