package repository

import (
	"context"
	"encoding/json"
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

func (gmr *GameMediaRepository) BatchInsert(ctx context.Context, gamesMedia []*models.GameMedia) error {
	conn, err := gmr.Pool.Acquire(ctx)
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
			INSERT INTO Games_media (appid, thumbnail_url, background_url, screenshots, movies)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (appid) DO UPDATE SET
				thumbnail_url = $2,
				background_url = $3,
				screenshots = $4,
				movies = $5
		`, gameMedia.AppID, gameMedia.ThumbnailURL, gameMedia.BackgroundURL, gameMedia.Screenshots, gameMedia.Movies)
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

func (gmr *GameMediaRepository) GetByAppID(ctx context.Context, appId int) (*models.GameMedia, error) {
	conn, err := gmr.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var gameMedia models.GameMedia
	var screenshotsJSON []byte
	var moviesJSON []byte

	err = conn.QueryRow(ctx, `
        SELECT appid, thumbnail_url, background_url, screenshots, movies
        FROM Games_media
        WHERE appid = $1
    `, appId).Scan(
		&gameMedia.AppID,
		&gameMedia.ThumbnailURL,
		&gameMedia.BackgroundURL,
		&screenshotsJSON,
		&moviesJSON,
	)
	if err != nil {
		return nil, err
	}

	if screenshotsJSON != nil {
		if err := json.Unmarshal(screenshotsJSON, &gameMedia.Screenshots); err != nil {
			return nil, fmt.Errorf("failed to unmarshal screenshots: %w", err)
		}
	} else {
		gameMedia.Screenshots = []models.Screenshot{} // Initialize empty slice
	}

	// Handle movies - check for null/empty before unmarshaling
	if moviesJSON != nil {
		if err := json.Unmarshal(moviesJSON, &gameMedia.Movies); err != nil {
			return nil, fmt.Errorf("failed to unmarshal movies: %w", err)
		}
	} else {
		gameMedia.Movies = []models.Movie{} // Initialize empty slice for games without movies
	}

	return &gameMedia, nil
}
