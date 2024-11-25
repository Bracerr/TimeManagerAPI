package database

import (
	"TimeManagerAuth/src/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"time"
)

func ConnectMongoDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbInitObject := config.GetDbParams()

	clientOptions := options.Client().ApplyURI(dbInitObject.Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	db := client.Database(dbInitObject.DbName)

	if err := migrateCollections(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateCollections(db *mongo.Database) error {
	if err := createNotionCollection(db); err != nil {
		return err
	}

	if err := createProjectCollection(db); err != nil {
		return err
	}

	if err := createReportCollection(db); err != nil {
		return err
	}

	if err := createUserCollection(db); err != nil {
		return err
	}

	return nil
}

func createNotionCollection(db *mongo.Database) error {
	err := db.CreateCollection(context.TODO(), "notions")
	return err
}

func createProjectCollection(db *mongo.Database) error {
	err := db.CreateCollection(context.TODO(), "projects")
	return err
}

func createReportCollection(db *mongo.Database) error {
	err := db.CreateCollection(context.TODO(), "reports")
	return err
}

func createUserCollection(db *mongo.Database) error {
	collection := db.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "login", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := collection.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		return err
	}
	return nil
}
