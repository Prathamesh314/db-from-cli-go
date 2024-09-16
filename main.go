package main

import (
	"bufio"
	"context"
	"db_cli/db_handler"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Available databases: ")
	fmt.Println("1. MongoDB")
	fmt.Println("2. MySQL")
	fmt.Println("3. Postgres")
	fmt.Println("Enter number of databse: ")
	var db_name string

	fmt.Scanln(&db_name)

	switch db_name {
	case "1":
		fmt.Println("Connecting to MongoDB...")
		var mongodb_connectin_uri string
		fmt.Println("Enter MongoDB connection URI: ")
		fmt.Scanln(&mongodb_connectin_uri)
		if mongodb_connectin_uri == "" {
			mongodb_connectin_uri = "mongodb+srv://affiliated:J42sd6P2aGwKZ1lN@cluster0.ufc2aiz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
		}
		var database_name string
		fmt.Println("Enter database name: ")
		fmt.Scanln(&database_name)
		var collection_name string
		fmt.Println("Enter collection name: ")
		fmt.Scanln(&collection_name)
		mongo_handler := db_handler.MongoHandler{
			MongoDBURI: mongodb_connectin_uri,
			DatabaseName: database_name,
			CollectionName: collection_name,
		}
		err := mongo_handler.ConnectToMongoDB()
		if err != nil {
			fmt.Println("Error connecting to MongoDB: ", err)
		}
		fmt.Println("Connected to MongoDB!")
		var option string
		fmt.Println("Choose what u want: ")
		fmt.Println("1. Change Collection")
		fmt.Println("2. Change Database")
		fmt.Println("3. Find One")
		fmt.Println("4. Find All")
		fmt.Println("5. Insert One")
		fmt.Println("6. Insert Many")
		fmt.Println("7. Update One")
		fmt.Println("8. Update Many")
		fmt.Println("9. Delete One")
		fmt.Println("10. Delete Many")
		fmt.Println("11. Exit")
		for {
				fmt.Print("> ")
			fmt.Scanln(&option)
			switch option {
			case "1":
				fmt.Println("Enter new collection name: ")
				fmt.Scanln(&collection_name)
			case "2":
				fmt.Println("Enter new database name: ")
				fmt.Scanln(&database_name)
			case "3":
				fmt.Println("Enter (key, value) to find: ")
				var key string
				var value string
				fmt.Scanln(&key, &value)
				fmt.Printf("Finding one document with %s: %s\n", key, value)
				result, err := mongo_handler.FindOne(key, value)
				if err != nil {
					fmt.Println("Error finding one document: ", err)
				}
				fmt.Println("Found one document: ", result)
			case "4":
				result, err := mongo_handler.FindAll()
				if err != nil {
					fmt.Println("Error finding all documents: ", err)
				}
				fmt.Println("Found all documents: ", result)
			}
			break
		}
	case "2":
		fmt.Println("Connecting to MySQL...")
	case "3":
		fmt.Println("Connecting to Postgres...")
		var postgres_user string
		fmt.Println("Enter Postgres user: ")
		fmt.Scanln(&postgres_user)
		var postgres_password string
		fmt.Println("Enter Postgres password: ")
		fmt.Scanln(&postgres_password)
		var postgres_host string
		fmt.Println("Enter Postgres host: ")
		fmt.Scanln(&postgres_host)
		var postgres_port string
		fmt.Println("Enter Postgres port: ")
		fmt.Scanln(&postgres_port)
		var postgres_database string
		fmt.Println("Enter Postgres database: ")
		fmt.Scanln(&postgres_database)
		postgres_handler := db_handler.PostgresHandler{
			POSTGRES_USER: postgres_user,
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

		// Scan the input
		scanner.Scan()
		query = scanner.Text()
		
		rows, err := postgres_handler.POOl.Query(context.Background(), query)
		if err != nil {
			fmt.Println("Error executing query: ", err)
			return
		}
		defer rows.Close()
		// Process the results
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				fmt.Println("Error scanning row: ", err)
				continue
			}
			
			// Print each row
			for _, value := range values {
				fmt.Printf("%v\n", value)
			}
			fmt.Println("---")
		}

		if err := rows.Err(); err != nil {
			fmt.Println("Error iterating rows: ", err)
		}

	default:
		fmt.Println("Invalid database name")
	}
}
