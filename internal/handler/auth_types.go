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

type RegisterRequest struct {
	Username string `json:"email" example:"buyer@test.com"`
	Password string `json:"password" example:"buyer123"`
	Role     string `json:"role" example:"buyer"`
}

type RegisterResponse struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"email" example:"buyer@test.com"`
	Role     string `json:"role" example:"buyer"`
}
