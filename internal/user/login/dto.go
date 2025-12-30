package login

type Output struct {
	Token string
}

type Input struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
