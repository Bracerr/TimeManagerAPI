package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Report struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Projects  []Project          `bson:"projects" json:"projects"`
	Context   string             `bson:"context" json:"context"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
