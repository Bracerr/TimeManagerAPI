package dto

type ProjectDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"  `
}

type UpdateProjectDto struct {
	Name           string `json:"name" validate:"required"`
	NewName        string `json:"newName"`
	NewDescription string `json:"newDescription"  `
}

type NotionDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"  `
	StartTime   string `json:"startTime" validate:"required"`
	EndTime     string `json:"endTime" validate:"required"`
	ProjectID   string `json:"projectID"  `
}

type DeleteNotionDto struct {
	ProjectID string `json:"projectID" validate:"required"`
	NotionID  string `json:"notionID" validate:"required"`
}

type UpdateNotionDto struct {
	ProjectID      string `json:"projectID" validate:"required"`
	NotionID       string `json:"notionID" validate:"required"`
	NewName        string `json:"newName"`
	NewDescription string `json:"newDescription"`
	NewStartTime   string `json:"newStartTime"`
	NewEndTime     string `json:"newEndTime"`
}
