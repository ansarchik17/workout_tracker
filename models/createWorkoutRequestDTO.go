package models

import "time"

type CreateWorkoutRequest struct {
	Title       string                         `json:"title" binding:"required"`
	Description string                         `json:"description" binding:"required"`
	Comment     string                         `json:"comment" binding:"required"`
	ScheduledAt time.Time                      `json:"scheduled_at" binding:"required"`
	Exercises   []CreateWorkoutExerciseRequest `json:"exercises" binding:"required,min=1"`
}
