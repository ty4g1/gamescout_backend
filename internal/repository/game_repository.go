package repository

import (
	"github.com/jackc/pgx/v5"
)

type GameRepository struct {
	conn *pgx.Conn
}

func NewGamesRepository(conn *pgx.Conn) *GameRepository {
	return &GameRepository{
		conn: conn,
	}
}
