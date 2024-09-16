package db_handler

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoHandler struct {
	DatabaseName string
	CollectionName string
	MongoDBURI string
	Collection *mongo.Collection
}


func NewMongoHandler(mongodb_connectin_uri string, collection_name string, database_name string) (*MongoHandler, error) {
	handler := &MongoHandler{
		DatabaseName: database_name,
		CollectionName: collection_name,
		MongoDBURI: mongodb_connectin_uri,
	}

	return handler, nil
}

func (self *MongoHandler)ConnectToMongoDB() error{
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


func (self *MongoHandler)FindAll() ([]bson.M, error) {
	cursor, err := self.Collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var results []bson.M
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
