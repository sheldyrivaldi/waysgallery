package middlewares

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFileProfile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		avatar, err := c.FormFile("avatar")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		avatarSrc, err := avatar.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		defer avatarSrc.Close()

		avatarTempFile, err := ioutil.TempFile("uploads", "image-*.png")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		defer avatarTempFile.Close()

		if _, err = io.Copy(avatarTempFile, avatarSrc); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		avatarData := avatarTempFile.Name()

		bannerFile, err := c.FormFile("banner")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		bannerSrc, err := bannerFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		defer bannerSrc.Close()

		bannerTempFile, err := ioutil.TempFile("uploads", "image-*.png")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		defer bannerTempFile.Close()

		if _, err = io.Copy(bannerTempFile, bannerSrc); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		bannerData := bannerTempFile.Name()

		c.Set("avatarFile", avatarData)
		c.Set("bannerFile", bannerData)
		return next(c)
	}
}
