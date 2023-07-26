package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	postdto "waysgallery/dto/post"
	resultdto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type handlerPost struct {
	PostRepositories repositories.PostRepositories
}
type dataPosts struct {
	Posts interface{} `json:"posts"`
}

func HandlerPost(PostRepositories repositories.PostRepositories) *handlerPost {
	return &handlerPost{PostRepositories}
}

func (h *handlerPost) FindPosts(c echo.Context) error {
	posts, err := h.PostRepositories.FindPosts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "errror 1"})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataPosts{
			Posts: posts,
		},
	})
}

func (h *handlerPost) GetPostByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	post, err := h.PostRepositories.GetPostByID(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataPosts{
			Posts: convertResponsePost(post),
		},
	})
}

func (h *handlerPost) CreatePost(c echo.Context) error {
	claims := c.Get("userLogin")
	id := claims.(jwt.MapClaims)["id"].(float64)
	userID := int(id)

	title := c.FormValue("title")
	description := c.FormValue("description")

	var postIsMatch = false
	var postID int
	for !postIsMatch {
		postID = int(time.Now().Unix())
		fmt.Println(postID)
		postData, _ := h.PostRepositories.GetPostByID(postID)
		fmt.Println(postData.ID)
		if postData.ID == 0 {
			postIsMatch = true
		}
	}

	newPost := models.Post{
		ID:          postID,
		Title:       title,
		Description: description,
		UserID:      userID,
	}

	post, err := h.PostRepositories.CreatePost(newPost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "error 2"})
	}
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	for i := 1; i <= 5; i++ {
		id := strconv.Itoa(i)
		filename := c.Get("imageFile" + id).(string)

		if filename != "" {
			resp, errUpload := cld.Upload.Upload(ctx, filename, uploader.UploadParams{Folder: "waysgallery"})
			if errUpload != nil {
				fmt.Println(errUpload.Error())
			}

			newPhoto := models.Photo{
				PostID: post.ID,
				URL:    resp.SecureURL,
			}

			_, err := h.PostRepositories.CreatePhoto(newPhoto)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "error 1"})
			}

			imageName := filename[8:]
			errRemove := os.Remove(fmt.Sprintf("uploads/%s", imageName))
			if errRemove != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": errRemove.Error()})
			}
		}
	}
	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataPosts{
			Posts: post,
		},
	})

}

func (h *handlerPost) FindPostByFollowingUser(c echo.Context) error {

	claims := c.Get("userLogin")
	id := claims.(jwt.MapClaims)["id"].(float64)
	currentUserID := int(id)

	user, err := h.PostRepositories.GetUserPostByID(currentUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "error 1"})
	}

	var posts []models.Post
	for _, following := range user.Followings {
		posts = append(posts, following.Post...)
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataPosts{
			Posts: posts,
		},
	})
}

func convertResponsePost(p models.Post) postdto.PostResponseDTO {
	return postdto.PostResponseDTO{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Photos:      p.Photos,
		User:        p.User,
	}
}
