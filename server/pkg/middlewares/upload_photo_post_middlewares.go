package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/labstack/echo/v4"
)

func UploadPhotoPost(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx = context.Background()
		var CLOUD_NAME = os.Getenv("CLOUD_NAME")
		var API_KEY = os.Getenv("API_KEY")
		var API_SECRET = os.Getenv("API_SECRET")
		cld, err := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		image1, err := c.FormFile("image1")
		if err != nil {
			c.Set("imagePost1", "")
		}
		src1, err := image1.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		defer src1.Close()
		resp1, errUpload := cld.Upload.Upload(ctx, src1, uploader.UploadParams{Folder: "waysgallery"})
		if errUpload != nil {
			return c.JSON(http.StatusInternalServerError, errUpload.Error())
		}
		fmt.Println(resp1.SecureURL)
		c.Set("imagePost1", resp1.SecureURL)

		image2, err := c.FormFile("image2")
		if err != nil {
			c.Set("imagePost2", "")
		}
		src2, err := image2.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		defer src2.Close()
		resp2, errUpload := cld.Upload.Upload(ctx, src2, uploader.UploadParams{Folder: "waysgallery"})
		if errUpload != nil {
			return c.JSON(http.StatusInternalServerError, errUpload.Error())
		}
		fmt.Println(resp2.SecureURL)
		c.Set("imagePost2", resp2.SecureURL)

		image3, err := c.FormFile("image3")
		if err != nil {
			c.Set("imagePost3", "")
		}
		src3, err := image3.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		defer src3.Close()
		resp3, errUpload := cld.Upload.Upload(ctx, src3, uploader.UploadParams{Folder: "waysgallery"})
		if errUpload != nil {
			return c.JSON(http.StatusInternalServerError, errUpload.Error())
		}
		fmt.Println(resp3.SecureURL)
		c.Set("imagePost3", resp3.SecureURL)

		image4, err := c.FormFile("image4")
		if err != nil {
			c.Set("imagePost4", "")
		}
		src4, err := image4.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		defer src4.Close()
		resp4, errUpload := cld.Upload.Upload(ctx, src4, uploader.UploadParams{Folder: "waysgallery"})
		if errUpload != nil {
			return c.JSON(http.StatusInternalServerError, errUpload.Error())
		}
		fmt.Println(resp4.SecureURL)
		c.Set("imagePost4", resp4.SecureURL)

		image5, err := c.FormFile("image5")
		if err != nil {
			c.Set("imagePost5", "")
		}
		src5, err := image5.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		defer src5.Close()
		resp5, errUpload := cld.Upload.Upload(ctx, src5, uploader.UploadParams{Folder: "waysgallery"})
		if errUpload != nil {
			return c.JSON(http.StatusInternalServerError, errUpload.Error())
		}
		fmt.Println(resp5.SecureURL)
		c.Set("imagePost5", resp5.SecureURL)
		return next(c)
	}
}
