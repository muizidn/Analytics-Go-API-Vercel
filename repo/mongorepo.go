package repo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoRepository(connectionString, database, collection string) (*MongoRepository, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &MongoRepository{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (repo *MongoRepository) Close() {
	if repo.client != nil {
		_ = repo.client.Disconnect(context.Background())
	}
}

func (repo *MongoRepository) Track(event string, properties map[string]interface{}) error {
	collection := repo.client.Database(repo.database).Collection(repo.collection)

	eventDocument := bson.M{
		"event":      event,
		"properties": properties,
		"timestamp":  time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), eventDocument)
	if err != nil {
		log.Printf("Error tracking event: %v", err)
		return err
	}

	return nil
}

func (repo *MongoRepository) FetchTracking(event string) ([]map[string]interface{}, error) {
	collection := repo.client.Database(repo.database).Collection(repo.collection)

	filter := bson.M{"event": event}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error fetching tracking events: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []map[string]interface{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Error decoding result: %v", err)
			continue
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error iterating over results: %v", err)
		return nil, err
	}

	return results, nil
}
