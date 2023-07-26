package routes

import (
	"waysgallery/handlers"
	"waysgallery/pkg/middlewares"
	"waysgallery/pkg/mysql"
	"waysgallery/repositories"

	"github.com/labstack/echo/v4"
)

func OrderRoutes(e *echo.Group) {
	orderRepository := repositories.RepositoryOrder(mysql.DB)
	h := handlers.HandlerOrder(orderRepository)

	e.GET("/orders", middlewares.Auth(h.FindOrders))
	e.GET("/orders/client/:clientID", middlewares.Auth(h.FindOrdersByClientID))
	e.GET("/orders/vendor/:vendorID", middlewares.Auth(h.FindOrdersByVendorID))
	e.GET("/order/:id", middlewares.Auth(h.GetOrderByID))
	e.POST("/order", middlewares.Auth(h.CreateOrder))
	e.PATCH("/order/:id", middlewares.Auth(h.UpdateOrderByID))
}
