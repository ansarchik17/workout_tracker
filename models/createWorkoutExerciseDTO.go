package models

type CreateWorkoutExerciseRequest struct {
	ExerciseID int     `json:"exercise_id" binding:"required"`
	Sets       int     `json:"sets" binding:"required,min=1"`
	Reps       int     `json:"reps" binding:"required,min=1"`
	Weight     float64 `json:"weight" binding:"required,gte=0"`
	Order      int     `json:"order"`
}
