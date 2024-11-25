package service

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/internal/dto"
	"TimeManagerAuth/src/internal/repository"
	"TimeManagerAuth/src/internal/scripts/primitiveConvert"
	"TimeManagerAuth/src/pkg/customErrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type NotionService struct {
	repo        *repository.NotionRepository
	projectRepo *repository.ProjectRepository
}

func NewNotionService(repo *repository.NotionRepository, projectRepo *repository.ProjectRepository) *NotionService {
	return &NotionService{repo: repo, projectRepo: projectRepo}
}

func (n *NotionService) CreateNotion(notion *domain.Notion, userID, projectID string) (*domain.Notion, error) {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return nil, err
	}

	primitiveProjectId, err := stringToObjectId(projectID)
	if err != nil {
		return nil, err
	}

	projectExist, err := n.projectRepo.ExistProjectByID(primitiveProjectId)
	if err != nil {
		return nil, err
	}

	if !projectExist {
		return nil, customErrors.NewAppError(http.StatusBadRequest, "project not found")
	}

	notion.Creator = domain.UserRef{
		Collection: "users",
		ID:         primitiveUserId,
	}

	notion.Project = domain.ProjectRef{
		Collection: "projects",
		ID:         primitiveProjectId,
	}

	resultNotion, err := n.repo.CreateNotion(notion)
	if err != nil {
		return nil, err
	}
	err = n.projectRepo.InsertNotionInProject(primitiveProjectId, *resultNotion)
	if err != nil {
		return nil, err
	}
	return resultNotion, nil
}

func (n *NotionService) GetUserNotions(userID string) ([]domain.Notion, error) {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return nil, err
	}
	return n.repo.FindNotionsByUserID(primitiveUserId)
}

func (n *NotionService) DeleteNotion(userID, notionID, projectID string) error {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return err
	}

	primitiveProjectId, err := stringToObjectId(projectID)
	if err != nil {
		return err
	}

	primitiveNotionId, err := stringToObjectId(notionID)
	if err != nil {
		return err
	}

	err = n.projectRepo.DeleteNotionFromProject(primitiveProjectId, primitiveNotionId)
	if err != nil {
		return err
	}
	return n.repo.DeleteNotionByUserID(primitiveUserId, primitiveNotionId)
}

func (n *NotionService) UpdateNotion(updateNotionDto *dto.UpdateNotionDto) (*domain.Notion, error) {
	primitiveNotionId, err := stringToObjectId(updateNotionDto.NotionID)
	if err != nil {
		return nil, err
	}
	notion, err := n.repo.FindNotionById(primitiveNotionId)
	if err != nil {
		return nil, err
	}

	if updateNotionDto.NewName != "" {
		notion.Name = updateNotionDto.NewName
	}

	if updateNotionDto.NewDescription != "" {
		notion.Description = updateNotionDto.NewDescription
	}

	if updateNotionDto.NewStartTime != "" {
		primitiveStartTime, err := primitiveConvert.StringToPrimitiveDate(updateNotionDto.NewStartTime)
		if err != nil {
			return nil, err
		}
		notion.StartTime = primitiveStartTime
	}

	if updateNotionDto.NewEndTime != "" {
		primitiveEndTime, err := primitiveConvert.StringToPrimitiveDate(updateNotionDto.NewEndTime)
		if err != nil {
			return nil, err
		}
		notion.EndTime = primitiveEndTime
	}

	notion.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	primitiveProjectId, err := stringToObjectId(updateNotionDto.ProjectID)
	if err != nil {
		return nil, err
	}

	err = n.projectRepo.UpdateNotionArray(notion, primitiveProjectId)
	if err != nil {
		return nil, err
	}

	return n.repo.UpdateNotion(notion)
}
