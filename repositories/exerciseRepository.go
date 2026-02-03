package repositories

import (
	"context"
	"workout-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseRepository struct {
	db *pgxpool.Pool
}

func NewExerciseRepository(conn *pgxpool.Pool) *ExerciseRepository {
	return &ExerciseRepository{db: conn}
}

func (repository *ExerciseRepository) Create(ctx context.Context, exercise models.Exercise) (int, error) {
	var exerciseID int

	err := repository.db.QueryRow(ctx,
		`insert into exercises (name, muscle_group)
		 values ($1, $2)
		 returning exercise_id`,
		exercise.Name,
		exercise.MuscleGroup,
	).Scan(&exerciseID)

	if err != nil {
		return 0, err
	}

	return exerciseID, nil
}
