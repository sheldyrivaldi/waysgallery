package main

import (
	"fmt"
	"net/http"
	"os"
	"waysgallery/database"
	"waysgallery/pkg/mysql"
	"waysgallery/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
	}))
	mysql.DatabaseInit()
	database.RunMigration()

	e.Static("/uploads", "./uploads")

	routes.RouteInit(e.Group("/api/v1"))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	fmt.Println("Server running on port 5000")
	e.Logger.Fatal(e.Start(":" + port))
}
