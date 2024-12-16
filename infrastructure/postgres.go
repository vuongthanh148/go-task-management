package infrastructure

import (
	"context"
	"fmt"
	"kn-assignment/internal/log"
	"kn-assignment/property"

	"github.com/exaring/otelpgx"
	"github.com/georgysavva/scany/dbscan"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context) (*pgxpool.Pool, *pgxscan.API) {
	pgCfg := property.Get().PostgresConfig
	pgx, scanapi, err := NewPostgresWithScanApi(ctx, pgCfg)
	if err != nil {
		panic(err)
	}

	return pgx, scanapi
}

func NewPostgresWithScanApi(ctx context.Context, pgCfg property.PostgresConfig) (*pgxpool.Pool, *pgxscan.API, error) {
	cfg, err := pgxpool.ParseConfig(pgCfg.ConnString)
	if err != nil {
		log.Fatalf(ctx, "unable to parse postgres connection uri: %v", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	// MaxConnLifetime
	if int64(pgCfg.MaxConnLifetime) > 0 {
		cfg.MaxConnLifetime = pgCfg.MaxConnLifetime
	}

	// MaxConnIdleTime
	if pgCfg.MaxConnIdleTime > 0 {
		cfg.MaxConnIdleTime = pgCfg.MaxConnIdleTime
	}

	// MaxConns
	if pgCfg.MaxConns > 0 {
		cfg.MaxConns = pgCfg.MaxConns
	}

	// MinConns
	if pgCfg.MinConns > 0 {
		cfg.MinConns = pgCfg.MinConns
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create pg connection: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("unable to ping db: %v", err)
	}

	scanApi := NewScanApi(ctx)
	return pool, scanApi, nil
}

func NewScanApi(ctx context.Context) *pgxscan.API {
	scanApi, err := pgxscan.NewDBScanAPI(dbscan.WithAllowUnknownColumns(true))
	if err != nil {
		log.Fatalf(ctx, "error create dbscanapi: %v", err)
	}

	api, err := pgxscan.NewAPI(scanApi)
	if err != nil {
		log.Fatalf(ctx, "error create sqlscanner api: %v", err)
	}
	return api
}
