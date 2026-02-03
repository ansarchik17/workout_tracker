package repositories

import (
	"context"
	"workout-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutRepository struct {
	db *pgxpool.Pool
}

func NewWorkoutRepository(conn *pgxpool.Pool) *WorkoutRepository {
	return &WorkoutRepository{db: conn}
}

func (repository *WorkoutRepository) Create(ctx context.Context, workout models.Workout) (int, error) {
	var workoutID int
	err := repository.db.QueryRow(
		ctx,
		`insert into workouts
		 (user_id, title, description, scheduled_at, status, comment)
		 values ($1, $2, $3, $4, $5, $6)
		 returning workout_id`,
		workout.UserID,
		workout.Title,
		workout.Description,
		workout.ScheduledAt,
		workout.Status,
		workout.Comment,
	).Scan(&workoutID)

	if err != nil {
		return 0, err
	}

	return workoutID, nil
}

func (repository *WorkoutRepository) AddExerciseToWorkout(ctx context.Context, exercise models.WorkoutExercise) (int, error) {
	var workoutExerciseID int

	err := repository.db.QueryRow(ctx,
		`insert into workout_exercises
		 (workout_id, exercise_id, sets, reps, weight, "order")
		 values ($1, $2, $3, $4, $5, $6)
		 returning id`,
		exercise.WorkoutID,
		exercise.ExerciseID,
		exercise.Sets,
		exercise.Reps,
		exercise.Weight,
		exercise.Order,
	).Scan(&workoutExerciseID)

	if err != nil {
		return 0, err
	}

	return workoutExerciseID, nil
}

func (r *WorkoutRepository) GetUserExercises(ctx context.Context, userID int) ([]models.Exercise, error) {
	rows, err := r.db.Query(
		ctx,
		`select distinct
  e.exercise_id,
  e.name,
  e.muscle_group
from exercises e
left join workout_exercises we on we.exercise_id = e.exercise_id
left join workouts w on w.workout_id = we.workout_id
where w.user_id = $1 or w.user_id is null
order by e.name;
`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise

	for rows.Next() {
		var e models.Exercise
		if err := rows.Scan(
			&e.ExerciseID,
			&e.Name,
			&e.MuscleGroup,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}

	return exercises, nil
}

func (repository *WorkoutRepository) GetUserWorkouts(ctx context.Context, userID int) ([]models.Workout, error) {

	rows, err := repository.db.Query(ctx,
		`select
			workout_id,
			user_id,
			title,
			description,
			comment,
			scheduled_at,
			status,
			created_at
		from workouts
		where user_id = $1
		order by scheduled_at desc`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []models.Workout

	for rows.Next() {
		var w models.Workout
		if err := rows.Scan(
			&w.WorkoutID,
			&w.UserID,
			&w.Title,
			&w.Description,
			&w.Comment,
			&w.ScheduledAt,
			&w.Status,
			&w.CreatedAt,
		); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}

	return workouts, nil
}
