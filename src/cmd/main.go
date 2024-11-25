package main

import (
	"TimeManagerAuth/src/internal/repository"
	"TimeManagerAuth/src/internal/router"
	"TimeManagerAuth/src/internal/service"
	"TimeManagerAuth/src/pkg/auth"
	"TimeManagerAuth/src/pkg/config"
	"TimeManagerAuth/src/pkg/database"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"log"
	"time"
)

// @title Time Manager API
// @version 1.0
// @description API documentation for Time Manager application
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	config.Init()
	e := echo.New()

	db, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.Client().Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	jwtParams := config.GetJwtParams()

	jwtManager := auth.NewManager(jwtParams)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)

	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)

	notionRepo := repository.NewNotionRepository(db)
	notionService := service.NewNotionService(notionRepo, projectRepo)

	e = echo.New()
	router.SetupRoute(e, userService, projectService, jwtManager, userRepo, notionService)

	e.Logger.Fatal(e.Start(":8080"))
}
