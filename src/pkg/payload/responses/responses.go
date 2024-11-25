package responses

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
