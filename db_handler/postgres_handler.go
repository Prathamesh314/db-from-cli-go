package db_handler

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"
)

type PostgresHandler struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_DATABASE string
	POOl              *pgxpool.Pool
}

func NewPostgresHandler(postgres_user, postgres_password, postgres_host, postgres_port, postgres_database string) (*PostgresHandler, error) {
	handler := &PostgresHandler{
		POSTGRES_USER:     postgres_user,
		POSTGRES_PASSWORD: postgres_password,
		POSTGRES_HOST:     postgres_host,
		POSTGRES_PORT:     postgres_port,
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

func (self *PostgresHandler) ShowHelp() {
	fmt.Println("0. List all Tables")
	fmt.Println("1. List Columns of a Table")
	fmt.Println("2. Execute a Query")
}

func (self *PostgresHandler) ListAllDatabases() error {
	query := "SELECT datname, pg_size_pretty(pg_database_size(datname)) AS size FROM pg_database WHERE datistemplate = false;"
	rows, err := self.POOl.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer rows.Close()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Database Name", "Size"})
	t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 30},
		{Number: 2, WidthMax: 10},
	})

	for rows.Next() {
		var dbName, dbSize string
		if err := rows.Scan(&dbName, &dbSize); err != nil {
			return err
		}
		t.AppendRow([]interface{}{dbName, dbSize})
	}

	if err := rows.Err(); err != nil {
		return err
	}

	fmt.Println("Databases in PostgreSQL:")
	t.Render()
	return nil
}

func (self *PostgresHandler) ListAllTables() error {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';"
	rows, err := self.POOl.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer rows.Close()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Table Name"})
	t.SetStyle(table.StyleLight)

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		t.AppendRow([]interface{}{tableName})
	}

	if err := rows.Err(); err != nil {
		return err
	}

	fmt.Println("Tables in PostgreSQL:")
	t.Render()
	return nil
}

func (self *PostgresHandler) ListColumnsOfTable(tableName string) error {
	query := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = '%s';", tableName)
	rows, err := self.POOl.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer rows.Close()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Column Name", "Data Type"})
	t.SetStyle(table.StyleLight)

	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return err
		}
		t.AppendRow([]interface{}{columnName, dataType})
	}

	if err := rows.Err(); err != nil {
		return err
	}

	fmt.Println("Columns in Table:", tableName)
	t.Render()
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..." // Add ellipsis to indicate truncation
	}
	return s
}

func (self *PostgresHandler) ExecuteQuery(query string) error {
	rows, err := self.POOl.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get column names
	columnNames := rows.FieldDescriptions()

	// Set up the table writer for output
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Limit each column width to handle horizontal overflow
	maxColumnWidth := 20 // Maximum column width (adjust as per your need)

	// Convert column names to a table.Row (which expects interface{} types)
	headers := make([]interface{}, len(columnNames))
	for i, col := range columnNames {
		headers[i] = truncateString(col.Name, maxColumnWidth)
	}
	t.AppendHeader(headers)
	t.SetStyle(table.StyleLight)

	// Iterate through rows and add each row to the table
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return err
		}

		// Convert the row values to a slice of interface{}, truncating if necessary
		rowValues := make([]interface{}, len(values))
		for i, value := range values {
			if value != nil {
				rowValues[i] = truncateString(fmt.Sprintf("%v", value), maxColumnWidth)
			} else {
				rowValues[i] = "NULL" // Handle NULL values
			}
		}

		t.AppendRow(rowValues)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Render the table to output the results
	fmt.Println("Query Results:")
	t.Render()
	return nil
}