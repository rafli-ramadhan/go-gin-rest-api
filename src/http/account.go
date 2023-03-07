package http

type GetUser struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"fullname"`
	Email          string `json:"email"`
	EmployeeNumber string `json:"employee_number"`
	Address        string `json:"address"`
	JobPosition    string `json:"job_position"`
	PhotoURL       string `json:"photo_url"`
}

type RegisterUser struct {
	Username  string `json:"username" validate:"required"`
	FullName  string `json:"fullname" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UpdateUser struct {
	Username       *string `json:"username"`
	FullName       *string `json:"fullname"`
	Email          *string `json:"email"`
	Password       *string `json:"password"`
	Address        *string `json:"address"`
	EmployeeNumber *string `json:"employee_number"`
	JobPosition    *string `json:"job_position"`
	KTPNumber      *int    `json:"ktp_number"`
	PhoneNumber    *string `json:"phone_number"`
	Gender         *string `json:"gender"`
	DOBString      *string `json:"date_of_birth" example:"yyyy-mm-dd"`
}
