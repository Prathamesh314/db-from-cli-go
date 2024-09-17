package main

import (
	"bufio"
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
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
	POSTGRES_DATABASE := os.Getenv("POSTGRES_DEFAULT_DB")

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
		var postgres_username string
		var postgres_password string
		var postgres_host string
		var postgres_port string
		var postgres_database string
		fmt.Println("Enter Postgres username: ")
		fmt.Scanln(&postgres_username)
		fmt.Println("Enter Postgres password: ")
		fmt.Scanln(&postgres_password)
		fmt.Println("Enter Postgres host: ")
		fmt.Scanln(&postgres_host)
		fmt.Println("Enter Postgres port: ")
		fmt.Scanln(&postgres_port)
		fmt.Println("Enter Postgres database: ")
		fmt.Scanln(&postgres_database)
		if postgres_username == "" {
			postgres_username = POSTGRES_USER
		}
		if postgres_password == "" {
			postgres_password = POSTGRES_PASSWORD
		}
		if postgres_host == "" {
			postgres_host = POSTGRES_HOST
		}
		if postgres_port == "" {
			postgres_port = POSTGRES_PORT
		}
		if postgres_database == "" {
			postgres_database = POSTGRES_DATABASE
		}
		postgres_handler := db_handler.PostgresHandler{
			POSTGRES_USER: postgres_username,
			POSTGRES_PASSWORD: postgres_password,
			POSTGRES_HOST: postgres_host,
			POSTGRES_PORT: postgres_port,
			POSTGRES_DATABASE: postgres_database,
		}
		fmt.Println(postgres_handler.POSTGRES_USER)
		fmt.Println(postgres_handler.POSTGRES_PASSWORD)
		fmt.Println(postgres_handler.POSTGRES_HOST)
		fmt.Println(postgres_handler.POSTGRES_PORT)
		fmt.Println(postgres_handler.POSTGRES_DATABASE)
		err := postgres_handler.ConnectToPostgres()
		if err != nil {
			fmt.Println("Error connecting to Postgres: ", err)
			return
		}
		fmt.Println("Connected to Postgres!")

		postgres_handler.ShowHelp()
		var command string
		for {
			fmt.Print("> ")
			fmt.Scanln(&command)
			switch command {
			case "0":
				err := postgres_handler.ListAllTables()
				if err != nil {
					fmt.Println("Error listing tables: ", err)
				}
			case "1":
				var tableName string
				fmt.Println("Enter table name: ")
				fmt.Scanln(&tableName)
				err := postgres_handler.ListColumnsOfTable(tableName)
				if err != nil {
					fmt.Println("Error listing columns: ", err)
				}
			case "2":
				scanner := bufio.NewScanner(os.Stdin)
				var query string
				fmt.Println("Enter query:")
				scanner.Scan()
				query = scanner.Text()
				fmt.Println("Query: ", query)
				err := postgres_handler.ExecuteQuery(query)
				if err != nil {
					fmt.Println("Error executing query: ", err)
				}
			case "help":
				postgres_handler.ShowHelp()
			case "clear":
				db_handler.ClearTerminal()
			case "exit":
				postgres_handler.Close()
				fmt.Println("Disconnected from Postgres!")
				return
			default:
				fmt.Println("Invalid command")
			}
		}

	default:
		fmt.Println("Invalid database name")
	}
}
