package routes

import (
	"waysgallery/handlers"
	"waysgallery/pkg/middlewares"
	"waysgallery/pkg/mysql"
	"waysgallery/repositories"

	"github.com/labstack/echo/v4"
)

func BankRoutes(e *echo.Group) {
	bankRepository := repositories.RepositoryBank(mysql.DB)
	h := handlers.HandlerBank(bankRepository)

	e.GET("/banks", h.FindBanks)
	e.POST("/bank", middlewares.Auth(h.CreateBank))
	e.GET("/bank/:id", h.GetBankByID)
	e.DELETE("/bank/:id", middlewares.Auth(h.DeleteBank))
}
