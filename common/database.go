package common

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool *pgxpool.Pool
)

func SetDatabase(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return errors.New("cannot assign database")
	}
	Pool = newPool
	return nil
}
