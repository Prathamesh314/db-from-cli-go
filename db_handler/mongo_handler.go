package db_handler

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoHandler struct {
	DatabaseName string
	CollectionName string
	MongoDBURI string
	Collection *mongo.Collection
	Client *mongo.Client
}


func NewMongoHandler(mongodb_connectin_uri string, collection_name string, database_name string) (*MongoHandler, error) {
	handler := &MongoHandler{
		DatabaseName: database_name,
		CollectionName: collection_name,
		MongoDBURI: mongodb_connectin_uri,
	}

	return handler, nil
}

func (self *MongoHandler) CloseConnection() error {
	err := self.Client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	return nil
}


func (self *MongoHandler) ChangeCollection(collection_name string) error {
	self.CollectionName = collection_name
	err := self.CloseConnection()
	if err != nil {
		return err
	}


	err = self.ConnectToMongoDB()
	if err != nil {
		return err
	}
	log.Default().Println("Collection changed to: ", collection_name)

	return nil
}

func (self *MongoHandler)ConnectToMongoDB() error{
	fmt.Println("Connecting to MongoDB...")
	clientOptions := options.Client().ApplyURI(self.MongoDBURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}


	fmt.Println("Connected to MongoDB!")
	if self.DatabaseName != "" && self.CollectionName != "" {
		collection := client.Database(self.DatabaseName).Collection(self.CollectionName)
		self.Collection = collection
	}
	self.Client = client

	return nil
}

func (self *MongoHandler)FindOne(key string, value string) (bson.M, error) {
	filter := bson.D{{Key: key, Value: value}}

	var result bson.M 
	err := self.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err 
	}

	return result, nil
}


func (self *MongoHandler)FindAll() ([]map[string]interface{}, error) {
	cursor, err := self.Collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (self *MongoHandler)InsertOne(data map[string]interface{}) error {
	insertion_data := bson.D{}
	for key, value := range data {
		insertion_data = append(insertion_data, bson.E{Key: key, Value: value})
	}

	_, err := self.Collection.InsertOne(context.TODO(), insertion_data)
	if err != nil {
		return err
	}

	return nil
}


func (self *MongoHandler) DeleteByID(id primitive.ObjectID) error {
    filter := bson.M{"_id": id}

    _, err := self.Collection.DeleteOne(context.TODO(), filter)
    if err != nil {
        return err
    }

    return nil
}

func (self *MongoHandler) DeleteOne(key string, value string) error {
	filter := bson.D{{Key: key, Value: value}}

	_, err := self.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (self *MongoHandler) DeleteAll() error {
	_, err := self.Collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}

	return nil
}

func (self *MongoHandler) ListAllDatabases() ([]string, error) {
    databases, err := self.Client.ListDatabaseNames(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    return databases, nil
}

func (self *MongoHandler) Help(){
	fmt.Println("0. List All Databases")
	fmt.Println("1. Change Collection")
	fmt.Println("2. Change Database")
	fmt.Println("3. Find One")
	fmt.Println("4. Find All")
	fmt.Println("5. Insert One")
	fmt.Println("6. Delete One")
	fmt.Println("7. Delete All")
	fmt.Println("8. Exit")
}

func (self *MongoHandler) MongoRunner(){
	for{
		var option string
	var collection_name string
	var database_name string

	fmt.Print("> ")
	fmt.Scanln(&option)
	switch option {
	case "0":
		err := self.ListDbAndCollections()
		if err != nil {
			fmt.Println("Error listing databases and collections: ", err)
		}
	case "1":
			fmt.Println("Enter new collection name: ")
			fmt.Scanln(&collection_name)
			err := self.ChangeCollection(collection_name)
			if err != nil {
				fmt.Println("Error changing collection: ", err)
			}

	case "2":
		fmt.Println("Enter new database name: ")
		fmt.Scanln(&database_name)
		fmt.Println("Enter new collection name: ")
		fmt.Scanln(&collection_name)
		err := self.ChangeDatabase(database_name, collection_name)
		if err != nil {
			fmt.Println("Error changing database: ", err)
		}
	case "3":
		fmt.Println("Enter (key, value) to find: ")
		var key string
		var value string
		fmt.Scanln(&key, &value)
		fmt.Printf("Finding one document with %s: %s\n", key, value)
		result, err := self.FindOne(key, value)
		if err != nil {
			fmt.Println("Error finding one document: ", err)
		}
		arr := []map[string]interface{}{result}
		PrettyPrint(arr);
	case "4":
		result, err := self.FindAll()
		if err != nil {
			fmt.Println("Error finding all documents: ", err)
		}
		PrettyPrint(result)
	case "5":
		document_map := make(map[string]interface{})
		for{
			var key, value string
			fmt.Println("Enter (key, value) to insert: ")
			fmt.Scanln(&key, &value)
			if key == "" || value == "" {
				break
			}
			document_map[key] = value
		}
		fmt.Println("Inserting document: ", document_map)
		err := self.InsertOne(document_map)
		if err != nil {
			fmt.Println("Error inserting one document: ", err)
		}
		fmt.Println("Document inserted successfully!")

	case "6":
		fmt.Println("Enter (key, value) to delete: ")
		var key, value string
		fmt.Scanln(&key, &value)
		if key == "_id" {
			// Handle _id as an ObjectID
			objectID, err := primitive.ObjectIDFromHex(value)
			if err != nil {
				fmt.Println("Error converting value to ObjectID: ", err)
				return
			}
			err = self.DeleteByID(objectID)
		} else {
			err := self.DeleteOne(key, value)
			if err != nil {
				fmt.Println("Error deleting one document: ", err)
				return
			}
		}
		fmt.Println("Document deleted successfully!")
	case "7":
		err := self.DeleteAll()
		if err != nil {
			fmt.Println("Error deleting all documents: ", err)
		}
		fmt.Println("All documents deleted successfully!")
	default:
		self.CloseConnection()
		fmt.Println("Disconnected from MongoDB!")
			return
		}
	}
}

func (self *MongoHandler) ShowDetails() {
	fmt.Println("Database Name: ", self.DatabaseName)
	fmt.Println("Collection Name: ", self.CollectionName)
}

func (self *MongoHandler) ChangeDatabase(database_name string, collection_name string) error {
	self.DatabaseName = database_name
	self.CollectionName = collection_name
	err := self.CloseConnection()
	if err != nil {
		return err
	}

	err = self.ConnectToMongoDB()
	if err != nil {
		return err
	}

	self.ShowDetails()
	return nil
}

func (self *MongoHandler) ListDbAndCollections() (error) {
	ClearTerminal()
	databases, err := self.ListAllDatabases()
	if err != nil {
		return err
	}
	for i, database := range databases {
		fmt.Println(i, ". ", database)
	}
	fmt.Println("\n")
	var database_index int
	fmt.Println("Enter database index: ")
	fmt.Scanln(&database_index)
	if !(database_index >= 0 && database_index<len(databases)) {
		fmt.Println("Invalid database index")
		return err
	}
	ClearTerminal()
	db := self.Client.Database(databases[database_index])
	collections, err := db.ListCollectionNames(context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}
	for i, collection := range collections {
		fmt.Println(i, ". ", collection)
	}
	fmt.Println("\n")
	var collection_index int
	fmt.Println("Enter collection index: ")
	fmt.Scanln(&collection_index)
	err = self.ChangeCollection(collections[collection_index])
	if err != nil {
		return err
	}
	if !(collection_index >= 0 && collection_index<len(collections)) {
		fmt.Println("Invalid collection index")
		return err
	}
	ClearTerminal()
	err = self.ChangeDatabase(databases[database_index], collections[collection_index])
	if err != nil {
		return err
	}
	return nil
}