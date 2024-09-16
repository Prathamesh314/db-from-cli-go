package db_handler

import (
	"context"
	"fmt"


	"github.com/jackc/pgx/v5/pgxpool"
)



func ConnectToPostgres(postgres_user, postgres_password, postgres_host, postgres_port, postgres_database string) (*pgxpool.Pool, error) {

    connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", postgres_user, postgres_password, postgres_host, postgres_port, postgres_database)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	
	return pool, nil
}
