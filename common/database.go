package common

import (
	"errors"
	"lexicon/lkpp-go-crawler/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool  *pgxpool.Pool
	Query *repository.Queries
)

func SetDatabase(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return errors.New("cannot assign database")
	}
	Pool = newPool
	return nil
}

func SetQuery(newQuery *repository.Queries) error {
	if newQuery == nil {
		return errors.New("cannot assign nil query")
	}
	Query = newQuery
	return nil
}
