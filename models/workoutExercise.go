package models

type WorkoutExercise struct {
	ID         int     `json:"id"`
	WorkoutID  int     `json:"workout_id"`
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Order      int     `json:"order"`
}
