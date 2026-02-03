package handlers

import (
	"net/http"
	"workout-tracker/models"
	"workout-tracker/repositories"

	"github.com/gin-gonic/gin"
)

type ExerciseHandlers struct {
	exerciseRepo *repositories.ExerciseRepository
}

func NewExerciseHandler(exerciseRepo *repositories.ExerciseRepository) *ExerciseHandlers {
	return &ExerciseHandlers{exerciseRepo: exerciseRepo}
}

func (handler *ExerciseHandlers) CreateExercise(c *gin.Context) {
	var req models.CreateExerciseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	exercise := models.Exercise{
		Name:        req.Name,
		MuscleGroup: req.MuscleGroup,
	}

	id, err := handler.exerciseRepo.Create(c.Request.Context(), exercise)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create exercise",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"exercise_id": id,
	})
}
