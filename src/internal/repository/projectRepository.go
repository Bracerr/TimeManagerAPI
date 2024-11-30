package repository

import (
	"TimeManagerAuth/src/internal/domain"
	"TimeManagerAuth/src/internal/dto"
	"TimeManagerAuth/src/pkg/customErrors"
	"TimeManagerAuth/src/pkg/payload/requests"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type ProjectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(db *mongo.Database) *ProjectRepository {
	return &ProjectRepository{
		collection: db.Collection("projects"),
	}
}

func removeNotionByID(notions []domain.Notion, idToRemove primitive.ObjectID) []domain.Notion {
	for i, notion := range notions {
		if notion.ID == idToRemove {
			notions[i] = notions[len(notions)-1]
			return notions[:len(notions)-1]
		}
	}
	return notions
}

func (r *ProjectRepository) CreateProject(project *domain.Project) (*domain.Project, error) {
	var existingProject domain.Project
	err := r.collection.FindOne(context.TODO(), bson.M{
		"name":        project.Name,
		"creator.$id": project.Creator.ID,
	}).Decode(&existingProject)

	if err == nil {
		return nil, customErrors.NewAppError(http.StatusConflict, "project already exists")
	}

	result, err := r.collection.InsertOne(context.TODO(), project)
	if err != nil {
		return nil, err
	}
	project.ID = result.InsertedID.(primitive.ObjectID)
	return project, nil
}

func (r *ProjectRepository) FindProjectsByUserID(creatorID primitive.ObjectID) ([]domain.Project, error) {
	filter := bson.M{"creator.$id": creatorID}
	var projects []domain.Project
	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) FindProjectByNameAndUserId(creatorID primitive.ObjectID,
	projectName string) (*domain.Project, error) {

	filter := bson.M{
		"creator.$id": creatorID,
		"name":        projectName,
	}
	var project domain.Project
	err := r.collection.FindOne(context.TODO(), filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, customErrors.NewAppError(http.StatusNotFound, "Project not found")
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindProjectById(projectId primitive.ObjectID) (*domain.Project, error) {
	filter := bson.M{
		"_id": projectId,
	}
	var project domain.Project
	err := r.collection.FindOne(context.TODO(), filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, customErrors.NewAppError(http.StatusNotFound, "Project not found")
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) DeleteProjectByUserID(creatorID primitive.ObjectID, projectName string) error {
	filter := bson.M{
		"creator.$id": creatorID,
		"name":        projectName,
	}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}
	if result.DeletedCount == 0 {
		return customErrors.NewAppError(http.StatusNotFound, "project not found")
	}
	return nil
}

func (r *ProjectRepository) UpdateProjectByUserID(creatorID primitive.ObjectID,
	updateProjectDto *dto.UpdateProjectDto) (*domain.Project, error) {
	filter := bson.M{
		"creator.$id": creatorID,
		"name":        updateProjectDto.Name,
	}

	project, err := r.FindProjectByNameAndUserId(creatorID, updateProjectDto.Name)
	if err != nil {
		return nil, err
	}

	if updateProjectDto.NewName != "" {
		project.Name = updateProjectDto.NewName
	}

	if updateProjectDto.NewDescription != "" {
		project.Description = updateProjectDto.NewDescription
	}

	project.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": project,
	}

	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, customErrors.NewAppError(http.StatusNotFound, "project not found or no changes made")
	}

	return project, nil
}

func (r *ProjectRepository) ExistProjectByID(projectID primitive.ObjectID) (bool, error) {
	var project domain.Project
	filter := bson.M{
		"_id": projectID,
	}

	err := r.collection.FindOne(context.TODO(), filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func (r *ProjectRepository) InsertNotionInProject(projectID primitive.ObjectID, notion domain.Notion) error {
	var project domain.Project
	filter := bson.M{
		"_id": projectID,
	}

	err := r.collection.FindOne(context.TODO(), filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customErrors.NewAppError(http.StatusNotFound, "Project not found")
		}
		return err
	}

	project.Notions = append(project.Notions, notion)
	project.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	update := bson.M{
		"$set": project,
	}
	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProjectRepository) DeleteNotionFromProject(projectID, notionID primitive.ObjectID) error {
	var project domain.Project
	filter := bson.M{
		"_id": projectID,
	}

	err := r.collection.FindOne(context.TODO(), filter).Decode(&project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return customErrors.NewAppError(http.StatusNotFound, "Project not found")
		}
		return err
	}

	project.Notions = removeNotionByID(project.Notions, notionID)
	project.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": project,
	}
	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProjectRepository) UpdateNotionArray(notion *domain.Notion, projectID primitive.ObjectID) error {
	filter := bson.M{
		"notions._id": notion.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"notions.$": notion,
		},
	}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return customErrors.NewAppError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (r *ProjectRepository) ProjectSearch(req *requests.ProjectSearchRequest) ([]domain.Project, error) {
	filter := bson.M{}

	if req.ProjectID != "" {
		primitiveProjectID, err := stringToObjectId(req.ProjectID)
		if err != nil {
			return nil, err
		}
		filter["project.$id"] = primitiveProjectID
	}

	if req.Name != "" {
		filter["name"] = bson.M{"$regex": req.Name, "$options": "i"}
	}

	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, customErrors.NewAppError(http.StatusInternalServerError, "ошибка получения записей")
	}
	defer cursor.Close(context.TODO())

	var projects []domain.Project
	if err = cursor.All(context.TODO(), &projects); err != nil {
		return nil, customErrors.NewAppError(http.StatusBadRequest, "ошибка преобразования записей")
	}
	return projects, nil
}
