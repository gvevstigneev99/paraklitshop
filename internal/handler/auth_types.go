package handler

type LoginRequest struct {
	Username string `json:"username" example:"buyer@test.com"`
	Password string `json:"password" example:"buyer123"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
