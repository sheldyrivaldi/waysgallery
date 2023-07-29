package models

type User struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	Fullname   string  `json:"fullname"`
	Email      string  `json:"email" gorm:"unique"`
	Password   string  `json:"-"`
	Avatar     string  `json:"avatar"`
	Greeting   string  `json:"greeting"`
	Banner     string  `json:"banner"`
	Balance    int     `json:"balance" gorm:"default:0"`
	Followings []*User `json:"followings" gorm:"many2many:user_followings"`
	Followers  []*User `json:"followers" gorm:"many2many:user_followers"`
	Post       []Post  `json:"post" gorm:"foreignKey:UserID"`
	Role       string  `json:"role" gorm:"default:'user'"`
}

type UserResponse struct {
	ID         int     `json:"id"`
	Email      string  `json:"email"`
	Fullname   string  `json:"fullname"`
	Role       string  `json:"role"`
	Post       []Post  `json:"post" gorm:"-"`
	Avatar     string  `json:"avatar"`
	Greeting   string  `json:"greeting"`
	Banner     string  `json:"banner"`
	Followings []*User `json:"followings" gorm:"many2many:user_followings"`
}

type UserResponseFollower struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
}

func (UserResponse) TableName() string {
	return "users"
}

func (UserResponseFollower) TableName() string {
	return "users"
}
