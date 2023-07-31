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
		var CONFIG_SENDER_NAME = "Landtick <sheldyrivaldi@gmail.com>"
		var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
		var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")
		var Amount = float64(withdraw.Amount)
		var AmountRupiah = rupiah.FormatRupiah(Amount)

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", CONFIG_SENDER_NAME)
		mailer.SetHeader("To", withdraw.User.Email)
		mailer.SetHeader("Subject", "Withdraw Status")
		mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
		<html style="-moz-osx-font-smoothing: grayscale; -webkit-font-smoothing: antialiased; background-color: #464646; margin: 0; padding: 0">
		  <head>
			<meta name="viewport" content="width=device-width, initial-scale=1" />
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
			<meta name="format-detection" content="telephone=no" />
			<title>GO Email Templates: Generic Template</title>
		  </head>
		  <body bgcolor="#d7d7d7" class="generic-template" style="-moz-osx-font-smoothing: grayscale; -webkit-font-smoothing: antialiased; background-color: #d7d7d7; margin: 0; padding: 0">
			<!-- Header Start -->
			<div class="bg-white header" bgcolor="#ffffff" style="background-color: white; width: 100%">
			  <table align="center" bgcolor="#ffffff" style="border-left: 10px solid white; border-right: 10px solid white; max-width: 600px; width: 100%">
				<tr height="80">
				  <td align="left" class="vertical-align-middle" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: middle">
					<a href="https://waysgallery-one.vercel.app/" target="_blank" style="-webkit-text-decoration-color: #f16522; color: #f16522; text-decoration: none; text-decoration-color: #f16522">
					  <img src="https://waysgallery-one.vercel.app/assets/logo-f646df01.svg" alt="GO" width="70" style="border: 0; font-size: 0; margin: 0; max-width: 100%; padding: 0" />
					</a>
				  </td>
				</tr>
			  </table>
			</div>
			<!-- Header End -->
		
			<!-- Content Start -->
			<table cellpadding="0" cellspacing="0" cols="1" bgcolor="#d7d7d7" align="center" style="max-width: 600px">
			  <tr bgcolor="#d7d7d7">
				<td height="50" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
			  </tr>
		
			  <!-- This encapsulation is required to ensure correct rendering on Windows 10 Mail app. -->
			  <tr bgcolor="#d7d7d7">
				<td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
				  <!-- Seperator Start -->
				  <table cellpadding="0" cellspacing="0" cols="1" bgcolor="#d7d7d7" align="center" style="max-width: 600px; width: 100%">
					<tr bgcolor="#d7d7d7">
					  <td height="30" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
				  </table>
				  <!-- Seperator End -->
		
				  <!-- Generic Pod Left Aligned with Price breakdown Start -->
				  <table align="center" cellpadding="0" cellspacing="0" cols="3" bgcolor="white" class="bordered-left-right" style="border-left: 10px solid #d7d7d7; border-right: 10px solid #d7d7d7; max-width: 600px; width: 100%">
					<tr height="50">
					  <td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr align="center">
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					  <td class="text-primary" style="color: #f16522; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
						<img src="http://dgtlmrktng.s3.amazonaws.com/go/emails/generic-email-template/tick.png" alt="GO" width="50" style="border: 0; font-size: 0; margin: 0; max-width: 100%; padding: 0" />
					  </td>
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr height="17">
					  <td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr align="center">
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					  <td class="text-primary" style="color: #f16522; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
						<h1 style="color: #f16522; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 30px; font-weight: 700; line-height: 34px; margin-bottom: 0; margin-top: 0">Withdraw received</h1>
					  </td>
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr height="30">
					  <td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr align="left">
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
						<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0">Hi %s,</p>
					  </td>
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr height="10">
					  <td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr align="left">
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
						<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0">Your withdraw was successful!</p>
						<br />
						<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0">
						  <strong>Payment Details:</strong><br />
		
						  Amount: %s <br />
						  Bank: %s <br />
						  Account: %d <br />
						</p>
						<br />
						<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0">We advise to keep this email for future reference.&nbsp;&nbsp;&nbsp;&nbsp;<br /></p>
					  </td>
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr height="30">
					  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					  <td style="border-bottom: 1px solid #d3d1d1; color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr height="30">
					  <td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
					<tr align="center">
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
						<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0"><strong>Withdraw ID: %d</strong></p>
						<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0">Status : %s</p>
						<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0"></p>
					  </td>
					  <td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
		
					<tr height="50">
					  <td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
				  </table>
				  <!-- Generic Pod Left Aligned with Price breakdown End -->
		
				  <!-- Seperator Start -->
				  <table cellpadding="0" cellspacing="0" cols="1" bgcolor="#d7d7d7" align="center" style="max-width: 600px; width: 100%">
					<tr bgcolor="#d7d7d7">
					  <td height="50" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
					</tr>
				  </table>
				  <!-- Seperator End -->
				</td>
			  </tr>
			</table>
			<!-- Content End -->
		
			<!-- Footer Start -->
			<div class="bg-gray-dark footer" bgcolor="#fff" height="165" style="background-color: #fff; width: 100%">
			  <table align="center" bgcolor="#fff" style="max-width: 600px; width: 100%">
				<tr height="15">
				  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
				</tr>
		
				<tr>
				  <td align="center" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
					<img src="https://waysgallery-one.vercel.app/assets/logo-f646df01.svg" alt="GO" width="50" style="border: 0; font-size: 0; margin: 0; max-width: 100%; padding: 0" />
				  </td>
				</tr>
		
				<tr height="2">
				  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
				</tr>
		
				<tr>
				  <td align="center" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top">
					<p class="text-black" style="color: black; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0">Copyright © Ways Gallery 2023. All rights reserved.</p>
					<p class="text-primary" style="color: #f16522; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0"></p>
				  </td>
				</tr>
		
				<tr height="15">
				  <td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top"></td>
				</tr>
			  </table>
			</div>
		  </body>
		</html>
		`, withdraw.User.Fullname, AmountRupiah, withdraw.Bank.Name, withdraw.AccountNumber, withdraw.ID, status))

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
