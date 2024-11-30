package handlers

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/internal/dto"
	"TimeManagerAuth/src/internal/service"
	"TimeManagerAuth/src/pkg/customErrors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type ProjectHandler struct {
	service   *service.ProjectService
	validator *validator.Validate
}

func NewProjectHandler(service *service.ProjectService, validator *validator.Validate) *ProjectHandler {
	return &ProjectHandler{service: service, validator: validator}
}

// CreateProject @Summary Create project
// @Description Создает проект на основе переданного JWT Token
// @Tags Projects
// @Accept json
// @Produce json
// @Param ProjectData body dto.ProjectDto true "Project data"
// @Success 201 {object} domain.Project
// @Failure 400 {object} customErrors.AppError "Wrong create form"
// @Failure 403 {object} customErrors.AppError "Invalid token"
// @Failure 409 {object} customErrors.AppError "project already exists"
// @Security ApiKeyAuth
// @Router /projects [post]
func (p *ProjectHandler) CreateProject(c echo.Context) error {
	projectDto := new(dto.ProjectDto)

	if err := c.Bind(projectDto); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, "ошибка преобразования данных в json")
	}

	if err := p.validator.Struct(projectDto); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	project := new(domain.Project)
	project.Name = projectDto.Name
	project.Description = projectDto.Description
	project.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	project.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	userID := c.Get("userID").(string)

	resultProject, err := p.service.CreateProject(project, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resultProject)
}

func (p *ProjectHandler) GetUsersProjects(c echo.Context) error {
	userID := c.Get("userID").(string)

	projects, err := p.service.GetUserProjects(userID)
	if err != nil {
		return err
	}
	if projects == nil {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, projects)
}

func (p *ProjectHandler) GetUserProject(c echo.Context) error {
	userID := c.Get("userID").(string)
	projectName := c.QueryParam("projectName")
	if projectName == "" {
		return customErrors.NewAppError(http.StatusBadRequest, "missing projectName")
	}
	project, err := p.service.GetUserProject(userID, projectName)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, project)
}

func (p *ProjectHandler) DeleteProject(c echo.Context) error {
	projectName := c.QueryParam("projectName")
	if projectName == "" {
		return customErrors.NewAppError(http.StatusBadRequest, "missing projectName")
	}
	userID := c.Get("userID").(string)

	err := p.service.DeleteProject(userID, projectName)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (p *ProjectHandler) UpdateProject(c echo.Context) error {
	userID := c.Get("userID").(string)
	updatedProjectDto := new(dto.UpdateProjectDto)

	if err := c.Bind(updatedProjectDto); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, "Ошибка преобразования данных в json")
	}

	if err := p.validator.Struct(updatedProjectDto); err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	project, err := p.service.UpdateProject(updatedProjectDto, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, project)
}

func (p *ProjectHandler) ProjectSearch(c echo.Context) error {
	projectId := c.QueryParam("id")

	//if err := c.Bind(projectId); err != nil {
	//	return customErrors.NewAppError(http.StatusInternalServerError, "Ошибка преобразования данных в json")
	//}
	//
	//if err := p.validator.Struct(projectId); err != nil {
	//	return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	//}

	projects, err := p.service.ProjectSearch(projectId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, projects)
}
