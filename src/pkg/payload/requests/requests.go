package requests

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpRequest struct {
	Login    string `bson:"login" json:"login" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
	Role     string `bson:"role" json:"role"`
}

type NotionSearchRequest struct {
	ProjectID   string `json:"projectId" validate:"required"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ProjectSearchRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
	Name      string `json:"name,omitempty"`
}
