package handlers

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/internal/dto"
	"TimeManagerAuth/src/internal/scripts/primitiveConvert"
	"TimeManagerAuth/src/internal/service"
	"TimeManagerAuth/src/pkg/customErrors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type NotionHandler struct {
	service   *service.NotionService
	validator *validator.Validate
}

func NewNotionHandler(service *service.NotionService, validator *validator.Validate) *NotionHandler {
	return &NotionHandler{service: service, validator: validator}
}

func (n *NotionHandler) CreateNotion(c echo.Context) error {
	notionDto := new(dto.NotionDto)

	if err := c.Bind(notionDto); err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, "ошибка преобразования данных в json")
	}

	if err := n.validator.Struct(notionDto); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	notion := new(domain.Notion)
	notion.Name = notionDto.Name

	if notionDto.Description != "" {
		notion.Description = notionDto.Description
	}

	notion.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	notion.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	var err error
	notion.StartTime, err = primitiveConvert.StringToPrimitiveDate(notionDto.StartTime)
	if err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	notion.EndTime, err = primitiveConvert.StringToPrimitiveDate(notionDto.EndTime)
	if err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("userID").(string)

	resultNotion, err := n.service.CreateNotion(notion, userID, notionDto.ProjectID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resultNotion)
}

func (n *NotionHandler) GetUsersNotions(c echo.Context) error {
	userID := c.Get("userID").(string)

	projects, err := n.service.GetUserNotions(userID)
	if err != nil {
		return err
	}
	if projects == nil {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, projects)
}

func (n *NotionHandler) DeleteNotion(c echo.Context) error {
	userID := c.Get("userID").(string)
	deleteNotionDto := new(dto.DeleteNotionDto)

	if err := c.Bind(deleteNotionDto); err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, "ошибка преобразования данных в json")
	}

	if err := n.validator.Struct(deleteNotionDto); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	err := n.service.DeleteNotion(userID, deleteNotionDto.NotionID, deleteNotionDto.ProjectID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (n *NotionHandler) UpdateNotion(c echo.Context) error {
	updateNotionDto := new(dto.UpdateNotionDto)
	if err := c.Bind(updateNotionDto); err != nil {
		return customErrors.NewAppError(http.StatusInternalServerError, "ошибка преобразования данных в json")
	}

	if err := n.validator.Struct(updateNotionDto); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	notion, err := n.service.UpdateNotion(updateNotionDto)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, notion)
}
