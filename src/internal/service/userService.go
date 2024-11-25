package service

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/internal/repository"
	"TimeManagerAuth/src/pkg/auth"
	"TimeManagerAuth/src/pkg/customErrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserService struct {
	repo       *repository.UserRepository
	jwtManager *auth.Manager
}

func NewUserService(repo *repository.UserRepository, jwtManager *auth.Manager) *UserService {
	return &UserService{repo: repo, jwtManager: jwtManager}
}

func (s *UserService) Signup(user *domain.User) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", customErrors.NewAppError(http.StatusInternalServerError, "error password -> hash")
	}

	if user.Role == "" {
		user.Role = "client"
	}

	if user.Role != "client" && user.Role != "office" && user.Role != "admin" {
		return "", customErrors.NewAppError(http.StatusBadRequest, "wrong role")
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err = s.repo.CreateUser(user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", customErrors.NewAppError(http.StatusConflict, "login not unique")
		}
		return "", customErrors.NewAppError(http.StatusInternalServerError, "internal user create error")
	}

	token, err := s.jwtManager.GenerateToken(user.Login, user.Role, user.ID)
	if err != nil {
		return "", customErrors.NewAppError(http.StatusInternalServerError, "error generating token")
	}

	return token, nil
}

func (s *UserService) Login(login, password string) (string, error) {
	user, err := s.repo.FindUserByLogin(login)
	if err != nil {
		return "", customErrors.NewAppError(http.StatusNotFound, "user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", customErrors.NewAppError(http.StatusForbidden, "wrong password")
	}
	token, err := s.jwtManager.GenerateToken(user.Login, user.Role, user.ID)
	if err != nil {
		return "", customErrors.NewAppError(http.StatusInternalServerError, "error generate token")
	}

	return token, nil
}
