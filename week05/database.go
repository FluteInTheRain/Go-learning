package main

import (
	"database/sql"
	"fmt"
)

func setupDB() (*sql.DB, error) {
	dsn := "postgres://postgres:password@localhost:5432/shop?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id    BIGSERIAL PRIMARY KEY,
		name  TEXT    NOT NULL,
		price NUMERIC NOT NULL
	);`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}
	return db, nil
}
