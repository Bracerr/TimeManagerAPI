package repository

import (
	"TimeManagerAuth/src/internal/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	result, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *UserRepository) FindUserByLogin(login string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.TODO(), bson.M{"login": login}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
