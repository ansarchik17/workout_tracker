package repositories

import (
	"context"
	"workout-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: conn}
}

func (r *AuthRepository) Create(ctx context.Context, user models.User) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		`insert into users (name, email, password_hash)
		 values ($1, $2, $3)
		 returning user_id`,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(&id)

	return id, err
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx,
		`select user_id, name, email, password_hash
		 from users
		 where email = $1`,
		email,
	).Scan(
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
