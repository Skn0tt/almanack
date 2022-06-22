package db

import (
	"context"
	"flag"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// AddFlags adds an option to the specified FlagSet that creates and tests a DB
func AddFlags(fl *flag.FlagSet, name, usage string) (q *Queries) {
	q = new(Queries)
	fl.Func(name, usage, func(dbURL string) error {
		q2, err := Open(dbURL)
		if q2 != nil {
			*q = *q2
		}
		return err
	})
	return
}

func Open(dbURL string) (q *Queries, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	db, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	q = New(db)
	return q, nil
}
