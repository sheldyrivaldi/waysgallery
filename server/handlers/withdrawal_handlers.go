package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	resultdto "waysgallery/dto/result"
	withdrawaldto "waysgallery/dto/withdrawal"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type handlerWithdrawal struct {
	WithdrawalRepositories repositories.WithdrawalRepositories
}
type dataWithdrawals struct {
	Withdrawals interface{} `json:"withdrawals"`
}

func HandlerWithdrawal(WithdrawalRepositories repositories.WithdrawalRepositories) *handlerWithdrawal {
	return &handlerWithdrawal{WithdrawalRepositories}
}

func (h *handlerWithdrawal) FindWithdrawals(c echo.Context) error {
	withdrawals, err := h.WithdrawalRepositories.FindWithdrawals()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataWithdrawals{
			Withdrawals: withdrawals,
		},
	})
}

func (h *handlerWithdrawal) FindWithdrawalsByUserID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	withdrawals, err := h.WithdrawalRepositories.FindWithdrawalsByUserID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataWithdrawals{
			Withdrawals: withdrawals,
		},
	})
}

func (h *handlerWithdrawal) CreateWithdrawal(c echo.Context) error {
	request := new(withdrawaldto.CreateWithdrawalDTO)

	errBind := c.Bind(&request)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: errBind.Error()})
	}

	validation := validator.New()
	validationErr := validation.Struct(request)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: validationErr.Error()})
	}

	var withdrawalIsMatch = false
	var withdrawalID int
	for !withdrawalIsMatch {
		withdrawalID = int(time.Now().Unix())
		fmt.Println(withdrawalID)
		withdrawalData, _ := h.WithdrawalRepositories.GetWithdrawalByID(withdrawalID)
		if withdrawalData.ID == 0 {
			withdrawalIsMatch = true
		}
	}

	claims := c.Get("userLogin")
	id := claims.(jwt.MapClaims)["id"].(float64)
	userID := int(id)

	parseBankID, _ := strconv.Atoi(request.BankID)
	parseAmount, _ := strconv.Atoi(request.Amount)
	parseAccountNumber, _ := strconv.Atoi(request.AccountNumber)

	user, err := h.WithdrawalRepositories.GetUserWithdrawalByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	if parseAmount > user.Balance {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Insufficient balance!"})
	}

	newWithdrawal := models.Withdrawal{
		ID:            withdrawalID,
		UserID:        userID,
		BankID:        parseBankID,
		AccountNumber: parseAccountNumber,
		Amount:        parseAmount,
	}

	withdrawal, err := h.WithdrawalRepositories.CreateWithdrawal(newWithdrawal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataWithdrawals{
			Withdrawals: convertResponseWithdrawal(withdrawal),
		},
	})
}

func (h *handlerWithdrawal) UpdateWithdrawal(c echo.Context) error {
	request := new(withdrawaldto.UpdateWithdrawalDTO)

	errBind := c.Bind(&request)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: errBind.Error()})
	}

	validation := validator.New()
	validationErr := validation.Struct(request)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: validationErr.Error()})
	}

	withdrawal_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	withdrawal, err := h.WithdrawalRepositories.GetWithdrawalByID(withdrawal_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}
	withdrawal.Status = request.Status

	data, err := h.WithdrawalRepositories.UpdateWithdrawal(withdrawal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	if data.Status == "Success" {
		user, errGetUser := h.WithdrawalRepositories.GetUserWithdrawalByID(data.UserID)
		if errGetUser != nil {
			return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errGetUser.Error()})
		}
		user.Balance = user.Balance - data.Amount
		_, errUpdateUser := h.WithdrawalRepositories.UpdateUserWithdrawal(user)
		if errUpdateUser != nil {
			return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errUpdateUser.Error()})
		}
	}
	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataWithdrawals{
			Withdrawals: convertResponseWithdrawal(data),
		},
	})

}

func convertResponseWithdrawal(w models.Withdrawal) withdrawaldto.WithdrawalResponseDTO {
	return withdrawaldto.WithdrawalResponseDTO{
		ID:            w.ID,
		User:          w.User,
		Bank:          w.Bank,
		AccountNumber: w.AccountNumber,
		Amount:        w.Amount,
		Status:        w.Status,
	}
}
