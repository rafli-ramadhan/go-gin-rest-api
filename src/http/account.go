package http

type GetUser struct {
	ID             string `json:"id" validate:"required"`
	Username       string `json:"username" validate:"required"`
	FullName       string `json:"fullname" validate:"required"`
	Email          string `json:"email" validate:"required"`
	JobPosition    string `json:"job_position" validate:"required"`
	EmployeeNumber string `json:"employee_number" validate:"required"`
	PhotoURL       string `json:"photo_url" validate:"required"`
}

type RegisterUser struct {
	Username  string `json:"username" validate:"required"`
	FullName  string `json:"fullname" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UpdateUser struct {
	Username       string  `json:"username"`
	FullName       string  `json:"fullname"`
	Email          string  `json:"email"`
	Address        string  `json:"address"`
	EmployeeNumber string  `json:"employee_number"`
	JobPosition    string  `json:"job_position"`
	KTPNumber      string  `json:"ktp_number"`
	PhoneNumber    string  `json:"phone_number"`
	Gender         string  `json:"gender"`
	DOBString      string  `json:"date_of_birth" example:"dd/mm/yyyy"`
}
