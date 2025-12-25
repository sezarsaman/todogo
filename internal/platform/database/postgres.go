package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	Conn *pgx.Conn
}

func New(ctx context.Context, dsn string) (*Postgres, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Postgres{Conn: conn}, nil
}
