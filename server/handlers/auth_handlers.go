package handlers

import (
	"net/http"
	"time"
	authdto "waysgallery/dto/auth"
	resultdto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/pkg/bcrypt"
	jwtToken "waysgallery/pkg/jwt"
	"waysgallery/repositories"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type handlerAuth struct {
	AuthRepositories repositories.AuthRepositories
}

type dataAuth struct {
	User interface{} `json:"user"`
}

func HandlerAuth(AuthRepositories repositories.AuthRepositories) *handlerAuth {
	return &handlerAuth{AuthRepositories}
}

func (h *handlerAuth) Login(c echo.Context) error {
	request := new(authdto.LoginRequestDTO)

	err := c.Bind(request)

	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	validation := validator.New()

	validationErr := validation.Struct(request)

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	user, err := h.AuthRepositories.Login(request.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	isValid := bcrypt.ComparePassword(request.Password, user.Password)
	if !isValid {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Username or Password Error!"})
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	token, err := jwtToken.GenerateToken(&claims)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	data := authdto.AuthResponseDTO{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Token:    token,
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataAuth{
			User: data,
		},
	})

}

func (h *handlerAuth) Register(c echo.Context) error {
	request := new(authdto.RegisterRequestDTO)

	err := c.Bind(request)

	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	validation := validator.New()

	validationErr := validation.Struct(request)

	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	isFound, _ := h.AuthRepositories.GetUserByEmail(request.Email)

	if isFound {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Email duplicated!"})
	}

	hashPassword, err := bcrypt.GeneratePassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	newUser := models.User{
		Fullname: request.Fullname,
		Email:    request.Email,
		Password: hashPassword,
	}

	data, err := h.AuthRepositories.Register(newUser)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataAuth{
			User: convertAuth(data),
		},
	})

}

func (h *handlerAuth) CheckAuth(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	user, _ := h.AuthRepositories.CheckAuth(int(userId))

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataAuth{
			User: convertAuth(user),
		},
	})
}

func convertAuth(u models.User) authdto.AuthRegisterResponseDTO {
	return authdto.AuthRegisterResponseDTO{
		ID:       u.ID,
		Fullname: u.Fullname,
		Email:    u.Email,
	}
}
