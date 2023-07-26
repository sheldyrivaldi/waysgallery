package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	resultdto "waysgallery/dto/result"
	userdto "waysgallery/dto/user"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type handlerUser struct {
	UserRepositories repositories.UserRepositories
}
type dataUsers struct {
	Users interface{} `json:"users"`
}

func HandlerUser(UserRepositories repositories.UserRepositories) *handlerUser {
	return &handlerUser{UserRepositories}
}

func (h *handlerUser) FindUsers(c echo.Context) error {
	users, err := h.UserRepositories.FindUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataUsers{
			Users: users,
		},
	})
}

func (h *handlerUser) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	user, err := h.UserRepositories.GetUserByID(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataUsers{
			Users: convertResponseUser(user),
		},
	})
}

func (h *handlerUser) UpdateUser(c echo.Context) error {
	request := new(userdto.UpdateUserDTO)

	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	fullname := c.FormValue("fullname")
	greeting := c.FormValue("greeting")
	avatar := c.Get("avatarFile").(string)
	banner := c.Get("bannerFile").(string)

	request.Fullname = fullname
	request.Greeting = greeting
	request.Avatar = avatar
	request.Banner = banner

	validation := validator.New()

	validationErr := validation.Struct(request)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: validationErr.Error()})
	}

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	resp, err := cld.Upload.Upload(ctx, avatar, uploader.UploadParams{Folder: "waysgallery"})
	if err != nil {
		fmt.Println(err.Error())
	}
	resp2, err := cld.Upload.Upload(ctx, banner, uploader.UploadParams{Folder: "waysgallery"})
	if err != nil {
		fmt.Println(err.Error())
	}

	user, err := h.UserRepositories.GetUserByID(user_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	if request.Fullname != "" {
		user.Fullname = request.Fullname
	}
	if request.Greeting != "" {
		user.Greeting = request.Greeting
	}
	if request.Avatar != "" {
		user.Avatar = resp.SecureURL
	}
	if request.Banner != "" {
		user.Banner = resp2.SecureURL
	}

	data, err := h.UserRepositories.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	avatarName := avatar[8:]
	errRemove := os.Remove(fmt.Sprintf("uploads/%s", avatarName))
	if errRemove != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errRemove.Error()})
	}
	bannerName := banner[8:]
	errRemove2 := os.Remove(fmt.Sprintf("uploads/%s", bannerName))
	if errRemove != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errRemove2.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataUsers{
			Users: convertResponseUser(data),
		},
	})
}

func (h *handlerUser) FollowingUser(c echo.Context) error {
	claims := c.Get("userLogin")
	id := claims.(jwt.MapClaims)["id"].(float64)
	currentUserID := int(id)

	request := new(userdto.FollowingUser)

	errBind := c.Bind(&request)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: errBind.Error()})
	}

	validation := validator.New()

	validationErr := validation.Struct(request)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: validationErr.Error()})
	}

	followingUserID, _ := strconv.Atoi(request.FollowingID)

	currentUser, err := h.UserRepositories.GetUserByID(currentUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	followingUser, err := h.UserRepositories.GetUserByID(followingUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	_, errFollowing := h.UserRepositories.FollowingUser(currentUser, followingUser)
	if errFollowing != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errFollowing.Error()})
	}

	_, errFollower := h.UserRepositories.FollowedByUser(currentUser, followingUser)
	if errFollower != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errFollower.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data:   "Success following user!",
	})

}

func convertResponseUser(u models.User) userdto.UserResponseDTO {
	return userdto.UserResponseDTO{
		ID:         u.ID,
		Fullname:   u.Fullname,
		Email:      u.Email,
		Avatar:     u.Avatar,
		Greeting:   u.Greeting,
		Banner:     u.Banner,
		Balance:    u.Balance,
		Followings: u.Followings,
		Followers:  u.Followers,
		Post:       u.Post,
		Role:       u.Role,
	}
}
