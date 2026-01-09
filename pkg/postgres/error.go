package postgres

import (
	"errors"

	"github.com/lib/pq"
)

func PostgresErr(err error) error {
	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return errors.New("this alias is already taken")
	}
	return err
}
