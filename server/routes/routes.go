package routes

import "github.com/labstack/echo/v4"

func RouteInit(e *echo.Group) {
	AuthRoutes(e)
	UserRoutes(e)
	PostRoutes(e)
	OrderRoutes(e)
	BankRoutes(e)
	WithdrawalRoutes(e)
	ProjectRoutes(e)
}
