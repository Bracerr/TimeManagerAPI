package handlers

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/internal/service"
	"TimeManagerAuth/src/pkg/customErrors"
	"TimeManagerAuth/src/pkg/payload/requests"
	"TimeManagerAuth/src/pkg/payload/responses"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	service   *service.UserService
	validator *validator.Validate
}

func NewUserHandler(service *service.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{service: service, validator: validator}
}

// Signup @Summary Signup a new user
// @Description Signup a new user and return a token
// @Tags users
// @Accept json
// @Produce json
// @Param signUpData body requests.SignUpRequest true "SignUp Data"
// @Success 201 {object} responses.TokenResponse
// @Failure 400 {object} customErrors.AppError "Wrong signUp form"
// @Failure 500 {object} customErrors.AppError "Internal error. Description in message"
// @Failure 409 {object} customErrors.AppError "Not unique login"
// @Router /users/signUp [post]
func (h *UserHandler) Signup(c echo.Context) error {
	signUpRequest := new(requests.SignUpRequest)
	if err := c.Bind(signUpRequest); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, "ошибка преобразования данных в json")
	}

	if err := h.validator.Struct(signUpRequest); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	user := domain.NewUser(*signUpRequest)
	token, err := h.service.Signup(user)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, responses.TokenResponse{Token: token})
}

// Login @Summary Login a user
// @Description Login a user and return a token
// @Tags users
// @Accept json
// @Produce json
// @Param loginData body requests.LoginRequest true "Login data"
// @Success 200 {object} responses.TokenResponse
// @Failure 400 {object} customErrors.AppError "Wrong login form"
// @Failure 404 {object} customErrors.AppError "User not found"
// @Failure 403 {object} customErrors.AppError "Bad credentials"
// @Router /users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	loginData := new(requests.LoginRequest)

	if err := c.Bind(loginData); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, "Ошибка преобразования данных в json")
	}

	if err := h.validator.Struct(loginData); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	token, err := h.service.Login(loginData.Login, loginData.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responses.TokenResponse{Token: token})
}
