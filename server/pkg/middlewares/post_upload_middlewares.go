package middlewares

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func UploadFilePost(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		for i := 1; i <= 5; i++ {
			id := strconv.Itoa(i)
			image, err := c.FormFile("image" + id)
			if err != nil {
				c.Set("imageFile"+id, "")
				continue
			}

			src, err := image.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			defer src.Close()

			tempFile, err := ioutil.TempFile("uploads", "image-*.png")
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			defer tempFile.Close()

			if _, err = io.Copy(tempFile, src); err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}

			data := tempFile.Name()

			c.Set("imageFile"+id, data)
		}

		return next(c)
	}
}
