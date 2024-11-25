package auth

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/pkg/config"
	"TimeManagerAuth/src/pkg/customErrors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type Manager struct {
	SigningKey     string
	ExpirationTime int
}

func NewManager(params config.JwtParams) *Manager {
	return &Manager{params.SigningKey, params.Expiration}
}

func (m *Manager) GenerateToken(login string, role string, id primitive.ObjectID) (string, error) {
	expirationDuration := time.Duration(m.ExpirationTime) * time.Hour
	expirationTime := time.Now().Add(expirationDuration).Unix()

	claims := domain.Claims{
		Login:      login,
		Role:       role,
		ID:         id,
		Expiration: expirationTime,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(m.SigningKey)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customErrors.NewAppError(http.StatusForbidden, "unexpected signing method")
		}
		return []byte(m.SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var id primitive.ObjectID
		if idStr, ok := claims["id"].(string); ok {
			id, err = primitive.ObjectIDFromHex(idStr)
			if err != nil {
				return nil, customErrors.NewAppError(http.StatusForbidden, "invalid ID format")
			}
		} else {
			return nil, customErrors.NewAppError(http.StatusForbidden, "missing ID in token claims")
		}

		domainClaims := &domain.Claims{
			Login:      claims["login"].(string),
			Role:       claims["role"].(string),
			ID:         id,
			Expiration: int64(claims["exp"].(float64)),
		}
		return domainClaims, nil
	}

	return nil, customErrors.NewAppError(http.StatusForbidden, "invalid token")
}
