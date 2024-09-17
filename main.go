package main

import (
	"bufio"
	"context"
	"db_cli/db_handler"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Available databases: ")
	fmt.Println("1. MongoDB")
	fmt.Println("2. MySQL")
	fmt.Println("3. Postgres")
	fmt.Println("Enter number of databse: ")
	var db_name string

	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("error loading .env file: %v", err)
	}

	mongodb_uri := os.Getenv("MONGODB_URI")
	POSTGRES_USER := os.Getenv("POSTGRES_USERNAME")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
	POSTGRES_DATABASE := os.Getenv("POSTGRES_DATABASE")

	fmt.Scanln(&db_name)

	switch db_name {
	case "1":
		fmt.Println("Connecting to MongoDB...")
		var mongodb_connectin_uri string
		fmt.Println("Enter MongoDB connection URI: ")
		fmt.Scanln(&mongodb_connectin_uri)
		if mongodb_connectin_uri == "" {
			mongodb_connectin_uri = mongodb_uri
		}
		mongo_handler := db_handler.MongoHandler{
			MongoDBURI: mongodb_connectin_uri,
		}
		err = mongo_handler.ConnectToMongoDB()
		if err != nil {
			fmt.Println("Error connecting to MongoDB: ", err)
			return
		}
		db_handler.ClearTerminal()
		mongo_handler.ListDbAndCollections()
		if err != nil {
			fmt.Println("Error connecting to MongoDB: ", err)
			return
		}
		fmt.Println("Connected to MongoDB!")
		db_handler.ClearTerminal()
		mongo_handler.ShowDetails()
		mongo_handler.Help()
		mongo_handler.MongoRunner()
	case "2":
		fmt.Println("Connecting to MySQL...")
	case "3":
		postgres_username := POSTGRES_USER
		postgres_password := POSTGRES_PASSWORD
		postgres_host := POSTGRES_HOST
		postgres_port := POSTGRES_PORT
		postgres_database := POSTGRES_DATABASE
		postgres_handler := db_handler.PostgresHandler{
			POSTGRES_USER: postgres_username,
			POSTGRES_PASSWORD: postgres_password,
			POSTGRES_HOST: postgres_host,
			POSTGRES_PORT: postgres_port,
			POSTGRES_DATABASE: postgres_database,
		}
		err := postgres_handler.ConnectToPostgres()
		if err != nil {
			fmt.Println("Error connecting to Postgres: ", err)
			return
		}
		fmt.Println("Connected to Postgres!")

		scanner := bufio.NewScanner(os.Stdin)
		var query string
		fmt.Println("Enter space-separated values:")

		scanner.Scan()
		query = scanner.Text()

		rows, err := postgres_handler.POOl.Query(context.Background(), query)
		if err != nil {
			fmt.Println("Error executing query: ", err)
			return
		}
		defer rows.Close()

		columnNames := rows.FieldDescriptions()
		results := []map[string]interface{}{}

		for rows.Next() {
			rowMap := make(map[string]interface{})

			values, err := rows.Values()
			if err != nil {
				fmt.Println("Error scanning row: ", err)
				continue
			}

			for i, value := range values {
				columnName := string(columnNames[i].Name)
				rowMap[columnName] = value
			}

			results = append(results, rowMap)
		}

		if err := rows.Err(); err != nil {
			fmt.Println("Error iterating rows: ", err)
		}

	default:
		fmt.Println("Invalid database name")
	}
}
