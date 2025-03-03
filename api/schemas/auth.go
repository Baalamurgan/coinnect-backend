package schemas

type SignupRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone        string `json:"phone" validate:"required"`
	AddressLine1 string `json:"address_line_1" validate:"required"`
	AddressLine2 string `json:"address_line_2"`
	AddressLine3 string `json:"address_line_3"`
	State        string `json:"state" validate:"required"`
	Pin          string `json:"pin" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line_1" validate:"required"`
	AddressLine2 string `json:"address_line_2"`
	AddressLine3 string `json:"address_line_3"`
	State        string `json:"state" validate:"required"`
	Pin          string `json:"pin" validate:"required"`
}
