package models

import "time"

type Exercise struct {
	ExerciseID  int       `json:"exercise_id"`
	Name        string    `json:"name"`
	MuscleGroup string    `json:"muscle_group"`
	CreatedAt   time.Time `json:"created_at"`
}
