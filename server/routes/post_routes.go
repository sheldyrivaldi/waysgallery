package routes

import (
	"waysgallery/handlers"
	"waysgallery/pkg/middlewares"
	"waysgallery/pkg/mysql"
	"waysgallery/repositories"

	"github.com/labstack/echo/v4"
)

func PostRoutes(e *echo.Group) {
	postRepository := repositories.RepositoryPost(mysql.DB)
	h := handlers.HandlerPost(postRepository)

	e.GET("/posts", h.FindPosts)
	e.GET("/posts/following", middlewares.Auth(h.FindPostByFollowingUser))
	e.POST("/post", middlewares.Auth(middlewares.UploadFilePost(h.CreatePost)))
	e.GET("/post/:id", h.GetPostByID)
}
