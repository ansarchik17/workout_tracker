package repositories

import (
	"context"
	"workout-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepositories(conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: conn}
}

func (repository *AuthRepository) Create(ctx context.Context, user models.User) (int, error) {
	var id int

	err := repository.db.QueryRow(ctx, "insert into users(name, email, password_hash) values ($1, $2, $3) returning id", user.Name, user.Email, user.PasswordHash).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, err
}
