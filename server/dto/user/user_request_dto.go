package userdto

type CreateUserDTO struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
type UpdateUserDTO struct {
	Fullname string `json:"fullname" form:"fullname"`
	Greeting string `json:"greeting" form:"greeting"`
	Avatar   string `json:"avatar" form:"avatar"`
	Banner   string `json:"banner" form:"banner"`
}

type FollowingUser struct {
	FollowingID string `json:"following_id" form:"following_id" validate:"required"`
}
