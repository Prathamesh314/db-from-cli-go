package db_handler
import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHandler struct {
	POSTGRES_USER string
	POSTGRES_PASSWORD string
	POSTGRES_HOST string
	POSTGRES_PORT string
	POSTGRES_DATABASE string
	POOl *pgxpool.Pool
}

func NewPostgresHandler(postgres_user, postgres_password, postgres_host, postgres_port, postgres_database string) (*PostgresHandler, error) {
	handler := &PostgresHandler{
		POSTGRES_USER: postgres_user,
		POSTGRES_PASSWORD: postgres_password,
		POSTGRES_HOST: postgres_host,
		POSTGRES_PORT: postgres_port,
		POSTGRES_DATABASE: postgres_database,
	}

	return handler, nil
}

func (self *PostgresHandler) ConnectToPostgres() error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", self.POSTGRES_USER, self.POSTGRES_PASSWORD, self.POSTGRES_HOST, self.POSTGRES_PORT, self.POSTGRES_DATABASE)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}

	self.POOl, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	return nil
}

func (self *PostgresHandler) Close() {
	self.POOl.Close()
}