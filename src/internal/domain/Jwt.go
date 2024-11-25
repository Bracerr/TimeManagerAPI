package domain

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	Login      string             `json:"login"`
	Role       string             `json:"role"`
	ID         primitive.ObjectID `json:"id"`
	Expiration int64              `json:"exp"`
	jwt.StandardClaims
}
