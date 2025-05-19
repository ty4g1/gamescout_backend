package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ty4g1/gamescout_backend/internal/models"
)

type GameRepository struct {
	Pool *pgxpool.Pool
}

func NewGamesRepository(pool *pgxpool.Pool) *GameRepository {
	return &GameRepository{
		Pool: pool,
	}
}

func (gr *GameRepository) BatchInsert(ctx context.Context, games []*models.Game) error {
	conn, err := gr.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, game := range games {
		batch.Queue(`
			INSERT INTO Games (appid, name, short_description, price, initial_price, discount, release_date, genres, tags, positive, negative, platforms, feature_vector)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			ON CONFLICT (appid) DO UPDATE SET
				name = $2,
                short_description = $3,
                price = $4,
                initial_price = $5,
                discount = $6,
                release_date = $7,
                genres = $8,
                tags = $9,
                positive = $10,
                negative = $11,
                platforms = $12,
                feature_vector = $13,
                last_updated = CURRENT_TIMESTAMP
		`, game.AppId, game.Name, game.ShortDesc, game.Price, game.InitialPrice, game.Discount, game.ReleaseDate, game.Genres, game.Tags, game.Positive, game.Negative, game.Platforms, game.FeatureVector)
	}

	br := tx.SendBatch(ctx, batch)

	for i := range batch.Len() {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("failed to insert the following game %v: %v", games[i], err)
		}
	}
	br.Close()

	if err := tx.Commit(ctx); err != nil {
		fmt.Println("yoooo")
		return err
	}

	return nil
}
