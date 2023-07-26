package handlers

import (
	"net/http"
	"strconv"
	"time"
	orderdto "waysgallery/dto/order"
	resultdto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
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

	orders, err := h.OrderRepositories.FindOrdersByClientID(vendorID)
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

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataOrders{
			Orders: order,
		},
	})

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
