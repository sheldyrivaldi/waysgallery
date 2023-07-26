package routes

import (
	"waysgallery/handlers"
	"waysgallery/pkg/middlewares"
	"waysgallery/pkg/mysql"
	"waysgallery/repositories"

	"github.com/labstack/echo/v4"
)

func WithdrawalRoutes(e *echo.Group) {
	withdrawalRepository := repositories.RepositoryWithdrawal(mysql.DB)
	h := handlers.HandlerWithdrawal(withdrawalRepository)

	e.GET("/withdrawals", middlewares.Auth(h.FindWithdrawals))
	e.POST("/withdrawal", middlewares.Auth(h.CreateWithdrawal))
	e.PATCH("/withdrawal/:id", middlewares.Auth(h.UpdateWithdrawal))
}
