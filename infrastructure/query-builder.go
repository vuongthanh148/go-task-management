package infrastructure

import "github.com/huandu/go-sqlbuilder"

func NewQueryBuilder() sqlbuilder.Flavor {
	return sqlbuilder.PostgreSQL
}
