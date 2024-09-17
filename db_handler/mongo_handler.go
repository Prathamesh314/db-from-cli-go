package db_handler

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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

	collection := client.Database(self.DatabaseName).Collection(self.CollectionName)
	self.Collection = collection
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

func (self *MongoHandler)InsertOne(data map[string]string) error {
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

func (self *MongoHandler) DeleteOne(key string, value string) error {
	filter := bson.D{{Key: key, Value: value}}

	_, err := self.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (self *MongoHandler) Help(){
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
