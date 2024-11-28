package router

import (
	"TimeManagerAuth/src/internal/customMiddleWare"
	"TimeManagerAuth/src/internal/handlers"
	"TimeManagerAuth/src/internal/repository"
	"TimeManagerAuth/src/internal/service"
	"TimeManagerAuth/src/pkg/auth"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func SetupRoute(e *echo.Echo, userService *service.UserService,
	projectService *service.ProjectService, jwtManager *auth.Manager, userRepository *repository.UserRepository, notionService *service.NotionService) {

	validate := validator.New()
	userHandler := handlers.NewUserHandler(userService, validate)
	projectHandler := handlers.NewProjectHandler(projectService, validate)
	notionHandler := handlers.NewNotionHandler(notionService, validate)

	userGroup := e.Group("/users", customMiddleWare.ErrorHandlerMiddleware)
	{
		userGroup.POST("/signUp", userHandler.Signup)
		userGroup.POST("/login", userHandler.Login)
	}

	projectGroup := e.Group("/projects", customMiddleWare.ErrorHandlerMiddleware, customMiddleWare.JWTMiddleware(jwtManager, userRepository))
	{
		projectGroup.POST("", projectHandler.CreateProject)
		projectGroup.GET("", projectHandler.GetUsersProjects)
		projectGroup.GET("/project", projectHandler.GetUserProject)
		projectGroup.DELETE("", projectHandler.DeleteProject)
		projectGroup.PUT("", projectHandler.UpdateProject)
	}

	notionGroup := e.Group("/notions", customMiddleWare.ErrorHandlerMiddleware, customMiddleWare.JWTMiddleware(jwtManager, userRepository))
	{
		notionGroup.POST("", notionHandler.CreateNotion)
		notionGroup.GET("", notionHandler.GetUsersNotions)
		notionGroup.DELETE("", notionHandler.DeleteNotion)
		notionGroup.PUT("", notionHandler.UpdateNotion)
		notionGroup.GET("/search", notionHandler.NotionSearch)
	}
}
