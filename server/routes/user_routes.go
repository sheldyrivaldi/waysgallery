package routes

import (
	"waysgallery/handlers"
	"waysgallery/pkg/middlewares"
	"waysgallery/pkg/mysql"
	"waysgallery/repositories"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerUser(userRepository)

	e.GET("/users", middlewares.Auth(h.FindUsers))
	e.POST("/user", middlewares.Auth(h.FollowingUser))
	e.GET("/user/:id", middlewares.Auth(h.GetUserByID))
	e.PATCH("/user/:id", middlewares.Auth(h.UpdateUser))
}
