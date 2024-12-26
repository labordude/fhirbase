package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	shared "github.com/labordude/fhirbase/shared"
)


// PgConnectionConfig holds PG credentials passed from command line
var PgConfig = shared.DatabaseOptions{}

func GetPgxConnectionConfig() (*pgxpool.Config, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		PgConfig.Username, PgConfig.Password, PgConfig.Host, PgConfig.Port, PgConfig.Database)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}
	config.MaxConns = 10
	config.MinConns = 2
	return config, nil
}

// GetConnection connects to database
func GetConnection() (*pgxpool.Pool, error) {
	conn, err := GetPgxConnectionConfig()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)

	if err != nil {
		log.Fatalf("Error testing database connection: %v", err)
	}

	return pool, nil
}
