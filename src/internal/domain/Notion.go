package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notion struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	StartTime   primitive.DateTime `bson:"startTime" json:"startTime"` // Формат: YYYY-MM-DDTHH:MM:SSZ
	EndTime     primitive.DateTime `bson:"endTime" json:"endTime"`     // Формат: YYYY-MM-DDTHH:MM:SSZ
	Project     ProjectRef         `bson:"project" json:"project"`
	Creator     UserRef            `bson:"user" json:"user"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

type ProjectRef struct {
	Collection string             `bson:"$ref" json:"$ref"`
	ID         primitive.ObjectID `bson:"$id" json:"$id"`
}
