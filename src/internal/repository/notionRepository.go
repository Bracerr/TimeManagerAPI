package repository

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/pkg/customErrors"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
)

type NotionRepository struct {
	collection *mongo.Collection
}

func stringToObjectId(id string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, customErrors.NewAppError(http.StatusInternalServerError, "ошибка преобразовани id из string -> primitive")
	}
	return objectId, nil
}

func NewNotionRepository(db *mongo.Database) *NotionRepository {
	return &NotionRepository{
		collection: db.Collection("notions"),
	}
}

func (r *NotionRepository) CreateNotion(notion *domain.Notion) (*domain.Notion, error) {
	var existingNotion domain.Notion
	err := r.collection.FindOne(context.TODO(), bson.M{
		"name":     notion.Name,
		"user.$id": notion.Creator.ID,
	}).Decode(&existingNotion)

	if err == nil {
		return nil, customErrors.NewAppError(http.StatusBadRequest, "notion already exists")
	}

	result, err := r.collection.InsertOne(context.TODO(), notion)
	if err != nil {
		return nil, err
	}
	notion.ID = result.InsertedID.(primitive.ObjectID)
	return notion, nil
}

func (r *NotionRepository) FindNotionsByUserID(creatorID primitive.ObjectID) ([]domain.Notion, error) {
	filter := bson.M{"user.$id": creatorID}
	var notions []domain.Notion
	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &notions)
	if err != nil {
		return nil, err
	}
	return notions, nil
}

func (r *NotionRepository) DeleteNotionByUserID(creatorID, notionID primitive.ObjectID) error {
	filter := bson.M{
		"user.$id": creatorID,
		"_id":      notionID,
	}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}
	if result.DeletedCount == 0 {
		return customErrors.NewAppError(http.StatusNotFound, "notion not found")
	}
	return nil
}

func (r *NotionRepository) FindNotionById(notionID primitive.ObjectID) (*domain.Notion, error) {
	filter := bson.M{
		"_id": notionID,
	}
	var notion domain.Notion
	err := r.collection.FindOne(context.TODO(), filter).Decode(&notion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, customErrors.NewAppError(http.StatusNotFound, "notion not found")
		}
		return nil, err
	}
	return &notion, nil
}

func (r *NotionRepository) UpdateNotion(notion *domain.Notion) (*domain.Notion, error) {
	filter := bson.M{
		"_id": notion.ID,
	}

	update := bson.M{
		"$set": notion,
	}

	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, customErrors.NewAppError(http.StatusNotFound, "notion not found or no changes made")
	}

	return notion, nil
}
