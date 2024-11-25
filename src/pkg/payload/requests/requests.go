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
