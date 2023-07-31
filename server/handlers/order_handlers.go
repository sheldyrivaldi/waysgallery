package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	orderdto "waysgallery/dto/order"
	resultdto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/yudapc/go-rupiah"
	"gopkg.in/gomail.v2"
)

type handlerOrder struct {
	OrderRepositories repositories.OrderRepositories
}
type dataOrders struct {
	Orders interface{} `json:"orders"`
}

func HandlerOrder(OrderRepositories repositories.OrderRepositories) *handlerOrder {
	return &handlerOrder{OrderRepositories}
}

func (h *handlerOrder) FindOrders(c echo.Context) error {
	orders, err := h.OrderRepositories.FindOrders()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataOrders{
			Orders: orders,
		},
	})
}

func (h *handlerOrder) FindOrdersByClientID(c echo.Context) error {

	clientID, _ := strconv.Atoi(c.Param("clientID"))

	orders, err := h.OrderRepositories.FindOrdersByClientID(clientID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataOrders{
			Orders: orders,
		},
	})
}
func (h *handlerOrder) FindOrdersByVendorID(c echo.Context) error {

	vendorID, _ := strconv.Atoi(c.Param("vendorID"))

	orders, err := h.OrderRepositories.FindOrdersByVendorID(vendorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataOrders{
			Orders: orders,
		},
	})
}

func (h *handlerOrder) CreateOrder(c echo.Context) error {
	request := new(orderdto.CreateOrderDTO)

	errBind := c.Bind(&request)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: errBind.Error()})
	}

	validation := validator.New()

	validationErr := validation.Struct(request)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: validationErr.Error()})
	}

	claims := c.Get("userLogin")
	id := claims.(jwt.MapClaims)["id"].(float64)
	orderByID := int(id)
	parsePrice, _ := strconv.Atoi(request.Price)
	parseOrderToID, _ := strconv.Atoi(request.OrderToID)

	if orderByID == parseOrderToID {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Cannot order to your account!"})

	}

	var orderIsMatch = false
	var orderID int
	for !orderIsMatch {
		orderID = int(time.Now().Unix())
		orderData, _ := h.OrderRepositories.GetOrderByID(orderID)
		if orderData.ID == 0 {
			orderIsMatch = true
		}
	}

	newOrder := models.Order{
		ID:          orderID,
		Title:       request.Title,
		Description: request.Description,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		Price:       parsePrice,
		OrderByID:   orderByID,
		OrderToID:   parseOrderToID,
	}

	order, err := h.OrderRepositories.CreateOrder(newOrder)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	user, err := h.OrderRepositories.GetUserOrderByID(order.OrderByID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(order.ID),
			GrossAmt: int64(order.Price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Fullname,
			Email: user.Email,
		},
	}

	snapResp, _ := s.CreateTransaction(req)

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data:   snapResp,
	})

}

func (h *handlerOrder) Notification(c echo.Context) error {
	var notificationPayload map[string]interface{}

	if err := c.Bind(&notificationPayload); err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed!", Message: err.Error()})
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderID := notificationPayload["order_id"].(string)

	order_id, _ := strconv.Atoi(orderID)

	data, err := h.OrderRepositories.GetOrderByID(order_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	fmt.Println("Ini payload", notificationPayload)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			_, err := h.OrderRepositories.UpdateOrderPayment("Waiting Accept", order_id)
			if err != nil {
				fmt.Println("Update order failed!")
			}
			SendMail("Waiting Accept", data)
		} else if fraudStatus == "accept" {
			_, err := h.OrderRepositories.UpdateOrderPayment("Waiting Accept", order_id)
			if err != nil {
				fmt.Println("Update order failed!")
			}
			SendMail("Waiting Accept", data)
		}
	} else if transactionStatus == "settlement" {
		_, err := h.OrderRepositories.UpdateOrderPayment("Waiting Accept", order_id)
		if err != nil {
			fmt.Println("Update order failed!")
		}
		SendMail("Waiting Accept", data)
	} else if transactionStatus == "deny" {
		_, err := h.OrderRepositories.UpdateOrderPayment("Cancel", order_id)
		if err != nil {
			fmt.Println("Update order failed!")
		}
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		_, err := h.OrderRepositories.UpdateOrderPayment("Cancel", order_id)
		if err != nil {
			fmt.Println("Update order failed!")
		}
	} else if transactionStatus == "pending" {
		SendMail("Waiting Accept", data)
		_, err := h.OrderRepositories.UpdateOrderPayment("Waiting Accept", order_id)
		if err != nil {
			fmt.Println("Update order failed!")
		}

	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{Status: "Success", Data: notificationPayload})
}

func (h *handlerOrder) GetOrderByID(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	order, err := h.OrderRepositories.GetOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataOrders{
			Orders: convertResponseOrder(order),
		},
	})
}

func (h *handlerOrder) UpdateOrderByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	request := new(orderdto.UpdateOrderDTO)

	errBind := c.Bind(&request)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: errBind.Error()})
	}

	validation := validator.New()

	validationErr := validation.Struct(request)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: validationErr.Error()})
	}

	order, err := h.OrderRepositories.GetOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	order.Status = request.Status

	data, err := h.OrderRepositories.UpdateOrder(order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	if data.Status == "Success" {
		user, err := h.OrderRepositories.GetUserOrderByID(data.OrderToID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
		}
		user.Balance = data.Price - (data.Price * 5 / 100)
		_, errUpdate := h.OrderRepositories.UpdateBalance(user)
		if errUpdate != nil {
			return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errUpdate.Error()})
		}
	}
	if data.Status == "Cancel" {
		user, err := h.OrderRepositories.GetUserOrderByID(data.OrderByID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
		}
		user.Balance = data.Price
		_, errUpdate := h.OrderRepositories.UpdateBalance(user)
		if errUpdate != nil {
			return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errUpdate.Error()})
		}
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataOrders{
			Orders: convertResponseOrder(data),
		},
	})

}

func SendMail(status string, order models.Order) {
	if status == "Waiting Accept" {
		var CONFIG_SMTP_HOST = "smtp.gmail.com"
		var CONFIG_SMTP_PORT = 587
		var CONFIG_SENDER_NAME = "Ways Gallery <sheldyrivaldi@gmail.com>"
		var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
		var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

		var TotalPayment = float64(order.Price)
		var TotalPaymentPupiah = rupiah.FormatRupiah(TotalPayment)

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", CONFIG_SENDER_NAME)
		mailer.SetHeader("To", order.OrderBy.Email)
		mailer.SetHeader("Subject", "Transaction Status")
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
						<li>Title : %s</li>
						<li>Description: %s</li>
						<li>StatusOrder : %s</li>
						<li>StatusPayment : %s</li>
						<li>Total payment: <b>%s</b></li>
					</ul>
				</body>
	  		</html>`, order.Title, order.Description, status, "Success", TotalPaymentPupiah))

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

		fmt.Println("Mail send to : " + order.OrderBy.Email)
	} else {
		fmt.Println("Error on sending mail!")
	}
}

func convertResponseOrder(m models.Order) orderdto.OrderResponseDTO {
	return orderdto.OrderResponseDTO{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		Price:       m.Price,
		OrderBy:     m.OrderBy,
		OrderTo:     m.OrderTo,
		Status:      m.Status,
	}
}
