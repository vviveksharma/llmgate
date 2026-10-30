package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	if dsn == "" {
		return &Postgres{}, nil
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	if p.DB == nil {
		return nil
	}
	return p.DB.PingContext(ctx)
}

func (p *Postgres) Close() error {
	if p.DB == nil {
		return nil
	}
	return p.DB.Close()
}
