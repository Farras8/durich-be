package requests

type RegisterAdmin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
