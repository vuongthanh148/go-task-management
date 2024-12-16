package userrepo

import (
	"kn-assignment/internal/core/port"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	dbPool     *pgxpool.Pool
	scanApi    *pgxscan.API
	sqlbuilder sqlbuilder.Flavor
}

func New(dbPool *pgxpool.Pool, scanApi *pgxscan.API, sqlbuilder sqlbuilder.Flavor) port.UserRepository {

	return &repository{
		dbPool:     dbPool,
		scanApi:    scanApi,
		sqlbuilder: sqlbuilder,
	}
}
