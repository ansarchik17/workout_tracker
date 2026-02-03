package handlers

import (
	"net/http"
	"time"
	"workout-tracker/models"
	"workout-tracker/repositories"

	"github.com/gin-gonic/gin"
)

type WorkoutHandler struct {
	workoutRepo *repositories.WorkoutRepository
}

func NewWorkoutHandler(workoutRepo *repositories.WorkoutRepository) *WorkoutHandler {
	return &WorkoutHandler{workoutRepo: workoutRepo}
}

func (handler *WorkoutHandler) CreateWorkout(c *gin.Context) {
	var req models.CreateWorkoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	userID := userIDValue.(int)

	workout := models.Workout{
		UserID:      userID,
		Title:       req.Title,
		Comment:     req.Comment,
		Description: req.Description,
		ScheduledAt: req.ScheduledAt,
		Status:      "planned",
		CreatedAt:   time.Now(),
	}

	workoutID, err := handler.workoutRepo.Create(c.Request.Context(), workout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create workout",
		})
		return
	}

	for _, ex := range req.Exercises {
		workoutExercise := models.WorkoutExercise{
			WorkoutID:  workoutID,
			ExerciseID: ex.ExerciseID,
			Sets:       ex.Sets,
			Reps:       ex.Reps,
			Weight:     ex.Weight,
			Order:      ex.Order,
		}

		_, err := handler.workoutRepo.AddExerciseToWorkout(
			c.Request.Context(),
			workoutExercise,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not add exercise to workout",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"workout_id": workoutID,
	})
}

func (handler *WorkoutHandler) GetMyExercises(c *gin.Context) {
	userID := c.GetInt("user_id")

	exercises, err := handler.workoutRepo.GetUserExercises(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not fetch exercises",
		})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

func (handler *WorkoutHandler) GetMyWorkouts(c *gin.Context) {
	userID := c.GetInt("user_id")

	workouts, err := handler.workoutRepo.GetUserWorkouts(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not fetch workouts",
		})
		return
	}

	c.JSON(http.StatusOK, workouts)
}
