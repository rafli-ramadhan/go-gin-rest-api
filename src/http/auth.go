package http

type Auth struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type AuthSelf struct {
	ID             string `json:"id" validate:"required"`
	Username       string `json:"username" validate:"required"`
	FullName       string `json:"fullname" validate:"required"`
	Email          string `json:"email" validate:"required"`
	JobPosition    string `json:"job_position" validate:"required"`
	EmployeeNumber string `json:"employee_number" validate:"required"`
	PhotoURL       string `json:"photo_url" validate:"required"`
}

type Token struct {
	Token string `json:"token"`
}

type ForgotPassword struct {
	Username    string `json:"username" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}