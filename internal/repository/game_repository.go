package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
		return err
	}

	return nil
}

func (gr *GameRepository) GetRandom(ctx context.Context, limit int, priceRange *models.PriceRange, releaseDate *models.ReleaseDate, tags []string, genres []string, platforms []string) ([]models.Game, error) {

	args := []any{limit, priceRange.Min, priceRange.Max}
	query := []string{`
		SELECT appid, name, short_description, price, initial_price, discount, 
           release_date, genres, tags, positive, negative, platforms, feature_vector
    FROM Games
    WHERE price BETWEEN $2 AND $3
	`}

	paramCount := 4

	if releaseDate != nil {
		dateOperator := ">"
		if releaseDate.IsBefore {
			dateOperator = "<"
		}
		args = append(args, releaseDate.Date)
		query = append(query, fmt.Sprintf("AND release_date %s $%d", dateOperator, paramCount))
		paramCount++
	}

	if tags != nil {
		args = append(args, tags)
		query = append(query, fmt.Sprintf("AND tags ?| $%d", paramCount))
		paramCount++
	}

	if genres != nil {
		args = append(args, genres)
		query = append(query, fmt.Sprintf("AND genres && $%d", paramCount))
		paramCount++
	}

	if platforms != nil {
		args = append(args, platforms)
		query = append(query, fmt.Sprintf("AND platforms && $%d", paramCount))
		paramCount++
	}

	query = append(query, `
		ORDER BY RANDOM()
		LIMIT $1
	`)

	conn, err := gr.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, strings.Join(query, "\n"), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.Game

	for rows.Next() {
		var game models.Game
		var tagsJSON []byte

		err := rows.Scan(
			&game.AppId,
			&game.Name,
			&game.ShortDesc,
			&game.Price,
			&game.InitialPrice,
			&game.Discount,
			&game.ReleaseDate,
			&game.Genres,
			&tagsJSON, // Scan JSONB as []byte first
			&game.Positive,
			&game.Negative,
			&game.Platforms,
			&game.FeatureVector,
		)
		if err != nil {
			return nil, err
		}

		// Parse the JSONB tags
		if err := json.Unmarshal(tagsJSON, &game.Tags); err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return games, nil
}
