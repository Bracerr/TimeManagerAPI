package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project @Description Project представляет собой структуру для проекта
// @Param CreatedAt query string true "Creation time in ISO 8601 format"
// @Param UpdatedAt query string true "Last updated time in ISO 8601 format"
type Project struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description"  `
	Creator     UserRef            `bson:"creator" json:"creator"`
	Notions     []Notion           `bson:"notions" json:"notions"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type UserRef struct {
	Collection string             `bson:"$ref" json:"$ref"`
	ID         primitive.ObjectID `bson:"$id" json:"$id"`
}
