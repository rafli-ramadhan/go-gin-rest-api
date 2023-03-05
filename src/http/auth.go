package http

type Auth struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type Token struct {
	Token string `json:"access_token"`
}

type ForgotPassword struct {
	KTPNumber int    `json:"ktp_number" validate:"required"`
	Password  string `json:"new_password" validate:"required"`
}
