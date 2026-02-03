package models

type CreateExerciseRequest struct {
	Name        string `json:"name" binding:"required"`
	MuscleGroup string `json:"muscle_group" binding:"required"`
}
