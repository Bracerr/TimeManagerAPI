package service

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/internal/dto"
	"TimeManagerAuth/src/internal/repository"
	"TimeManagerAuth/src/pkg/customErrors"
	"TimeManagerAuth/src/pkg/payload/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func stringToObjectId(id string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, customErrors.NewAppError(http.StatusInternalServerError, "ошибка преобразовани id из string -> primitive")
	}
	return objectId, nil
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (p *ProjectService) CreateProject(project *domain.Project, userID string) (*domain.Project, error) {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return nil, err
	}
	project.Creator = domain.UserRef{
		Collection: "users",
		ID:         primitiveUserId,
	}

	resultProject, err := p.repo.CreateProject(project)
	if err != nil {
		return nil, err
	}

	return resultProject, nil
}

func (p *ProjectService) GetUserProjects(userID string) ([]domain.Project, error) {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return nil, err
	}
	return p.repo.FindProjectsByUserID(primitiveUserId)
}

func (p *ProjectService) GetUserProject(userID, projectName string) (*domain.Project, error) {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return nil, err
	}
	return p.repo.FindProjectByNameAndUserId(primitiveUserId, projectName)
}

func (p *ProjectService) DeleteProject(userID, projectName string) error {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return err
	}

	return p.repo.DeleteProjectByUserID(primitiveUserId, projectName)
}

func (p *ProjectService) UpdateProject(updateProjectDto *dto.UpdateProjectDto, userID string) (*domain.Project, error) {
	primitiveUserId, err := stringToObjectId(userID)
	if err != nil {
		return nil, err
	}

	project, err := p.repo.UpdateProjectByUserID(primitiveUserId, updateProjectDto)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (p *ProjectService) ProjectSearch(req *requests.ProjectSearchRequest) ([]domain.Project, error) {
	projects, err := p.repo.ProjectSearch(req)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
