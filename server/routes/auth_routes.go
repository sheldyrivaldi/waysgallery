package routes

import (
	"waysgallery/handlers"
	"waysgallery/pkg/middlewares"
	"waysgallery/pkg/mysql"
	"waysgallery/repositories"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Group) {
	authRepository := repositories.RepositoryAuth(mysql.DB)
	h := handlers.HandlerAuth(authRepository)

	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.GET("/check-auth", middlewares.Auth(h.CheckAuth))
}
