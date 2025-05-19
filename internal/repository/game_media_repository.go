package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ty4g1/gamescout_backend/internal/models"
)

type GameMediaRepository struct {
	Pool *pgxpool.Pool
}

func NewGameMediaRepository(pool *pgxpool.Pool) *GameMediaRepository {
	return &GameMediaRepository{
		Pool: pool,
	}
}

func (gr *GameMediaRepository) BatchInsert(ctx context.Context, gamesMedia []*models.GameMedia) error {
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

	for _, gameMedia := range gamesMedia {
		batch.Queue(`
			INSERT INTO Games_media (appid, thumbnail_url, background_url, screenshots)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (appid) DO UPDATE SET
				thumbnail_url = $2,
				background_url = $3,
				screenshots = $4
		`, gameMedia.AppID, gameMedia.ThumbnailURL, gameMedia.BackgroundURL, gameMedia.Screenshots)
	}

	br := tx.SendBatch(ctx, batch)

	for i := range batch.Len() {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("failed to insert the following game %v: %v", gamesMedia[i], err)
		}
	}
	br.Close()

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
