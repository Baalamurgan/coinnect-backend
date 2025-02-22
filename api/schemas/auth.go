package schemas

type SignupRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
