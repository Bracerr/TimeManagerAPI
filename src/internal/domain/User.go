package domain

import (
	"TimeManagerAuth/src/pkg/payload/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Login     string             `bson:"login" json:"login" validate:"required"`
	Password  string             `bson:"password" json:"password" validate:"required"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

func NewUser(request requests.SignUpRequest) *User {
	return &User{
		Login:    request.Login,
		Password: request.Password,
		Role:     request.Role,
	}
}
