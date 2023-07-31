package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	resultdto "waysgallery/dto/result"
	withdrawaldto "waysgallery/dto/withdrawal"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/yudapc/go-rupiah"
	"gopkg.in/gomail.v2"
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

		SendMailWithdraw("Success", withdrawal)
	}
	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataWithdrawals{
			Withdrawals: convertResponseWithdrawal(data),
		},
	})

}

func SendMailWithdraw(status string, withdraw models.Withdrawal) {
	if status == "Success" {
		var CONFIG_SMTP_HOST = "smtp.gmail.com"
		var CONFIG_SMTP_PORT = 587
		var CONFIG_SENDER_NAME = "Ways Gallery <sheldyrivaldi@gmail.com>"
		var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
		var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")
		var Amount = float64(withdraw.Amount)
		var AmountRupiah = rupiah.FormatRupiah(Amount)

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", CONFIG_SENDER_NAME)
		mailer.SetHeader("To", withdraw.User.Email)
		mailer.SetHeader("Subject", "Withdraw Status")
		mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
			<html lang="en">
				<head>
					<meta charset="UTF-8" />
					<meta http-equiv="X-UA-Compatible" content="IE=edge" />
					<meta name="viewport" content="width=device-width, initial-scale=1.0" />
					<title>Document</title>
					<style>
						h1 {
						color: brown;
						}
					</style>
				</head>
				<body>
					<h2>Ticket Payment :</h2>
					<ul style="list-style-type:none;">
						<li>Withdrawal ID : %d</li>
						<li>Bank: %s</li>
						<li>Account Number: %d</li>
						<li>Status: %s</li>
						<li>Amount : <b>%s</b></li>
					</ul>
				</body>
	  		</html>`, withdraw.ID, withdraw.Bank.Name, withdraw.AccountNumber, status, AmountRupiah))

		dialer := gomail.NewDialer(
			CONFIG_SMTP_HOST,
			CONFIG_SMTP_PORT,
			CONFIG_AUTH_EMAIL,
			CONFIG_AUTH_PASSWORD,
		)

		err := dialer.DialAndSend(mailer)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Mail send to : " + withdraw.User.Email)
	} else {
		fmt.Println("Error on sending mail!")
	}
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
