package routes

import (
	"waysgallery/handlers"
	"waysgallery/pkg/middlewares"
	"waysgallery/pkg/mysql"
	"waysgallery/repositories"

	"github.com/labstack/echo/v4"
)

func ProjectRoutes(e *echo.Group) {
	projectRepository := repositories.RepositoryProject(mysql.DB)
	h := handlers.HandlerProject(projectRepository)

	e.POST("/project", middlewares.Auth(middlewares.UploadFileProject(h.CreateProject)))
	e.GET("/project/:id", middlewares.Auth(h.GetProjectByID))
	e.GET("/project/order/:id", middlewares.Auth(h.GetProjectByOrderID))
	e.PATCH("/project/:id", middlewares.Auth(middlewares.UploadFileProject(h.UpdateProject)))
}
