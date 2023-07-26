package authdto

type AuthResponseDTO struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
type AuthRegisterResponseDTO struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}
