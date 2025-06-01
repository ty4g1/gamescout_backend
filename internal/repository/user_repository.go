package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ty4g1/gamescout_backend/internal/models"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Pool: pool,
	}
}

func (ur *UserRepository) AddUser(ctx context.Context, id string) (*models.User, error) {
	conn, err := ur.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var user models.User

	err = conn.QueryRow(ctx, `
		INSERT INTO Users (cookie_id, swipe_history, preference_vector)
		VALUES ($1, $2, $3)
		ON CONFLICT (cookie_id) DO UPDATE SET
			swipe_history = $2,
			preference_vector = $3
		RETURNING cookie_id, swipe_history, preference_vector
	`, id, []string{}, make([]float64, 300)).Scan(&user.ID, &user.SwipeHistory, &user.PreferenceVector)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserPreference(ctx context.Context, id string) ([]float64, error) {
	conn, err := ur.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var preferenceVector []float64

	err = conn.QueryRow(ctx, `
		SELECT preference_vector FROM Users
		WHERE cookie_id = $1
	`, id).Scan(&preferenceVector)
	if err != nil {
		return nil, err
	}
	return preferenceVector, nil
}

func (ur *UserRepository) GetUserSwipes(ctx context.Context, id string) ([]int, error) {
	conn, err := ur.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var swipeHistory []int

	err = conn.QueryRow(ctx, `
		SELECT swipe_history FROM Users
		WHERE cookie_id = $1
	`, id).Scan(&swipeHistory)
	if err != nil {
		return nil, err
	}
	return swipeHistory, nil
}

func (ur *UserRepository) UpdateUserPreference(ctx context.Context, id string, preferenceVector []float64) error {
	conn, err := ur.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		UPDATE Users
		SET preference_vector = $1
		WHERE cookie_id = $2
	`, preferenceVector, id)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateUserSwipes(ctx context.Context, id string, swipes []int) error {
	conn, err := ur.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		UPDATE Users
		SET swipe_history = swipe_history || $1
		WHERE cookie_id = $2
	`, swipes, id)

	if err != nil {
		return err
	}

	return nil
}
