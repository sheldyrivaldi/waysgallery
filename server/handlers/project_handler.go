package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	projectdto "waysgallery/dto/project"
	resultdto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/labstack/echo/v4"
)

type handlerProject struct {
	ProjectRepositories repositories.ProjectRepositories
}
type dataProject struct {
	Project interface{} `json:"Project"`
}

func HandlerProject(ProjectRepositories repositories.ProjectRepositories) *handlerProject {
	return &handlerProject{ProjectRepositories}
}

func (h *handlerProject) CreateProject(c echo.Context) error {

	description := c.FormValue("description")
	orderID, _ := strconv.Atoi(c.FormValue("order_id"))

	var projectIsMatch = false
	var projectID int
	for !projectIsMatch {
		projectID = int(time.Now().Unix())
		projectData, _ := h.ProjectRepositories.GetProjectByID(projectID)
		if projectData.ID == 0 {
			projectIsMatch = true
		}
	}

	newProject := models.Project{
		ID:          projectID,
		Description: description,
		OrderID:     orderID,
	}

	project, err := h.ProjectRepositories.CreateProject(newProject)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	var PhotosProject []string

	for i := 1; i <= 5; i++ {
		id := strconv.Itoa(i)
		image, err := c.FormFile("image" + id)
		if err == nil {
			src, err := image.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			defer src.Close()

			var ctx = context.Background()
			var CLOUD_NAME = os.Getenv("CLOUD_NAME")
			var API_KEY = os.Getenv("API_KEY")
			var API_SECRET = os.Getenv("API_SECRET")
			cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

			resp, errUpload := cld.Upload.Upload(ctx, src, uploader.UploadParams{Folder: "waysgallery"})
			if errUpload != nil {
				fmt.Println(errUpload.Error())
			}

			PhotosProject = append(PhotosProject, resp.SecureURL)
		}
	}

	if len(PhotosProject) != 0 {
		for _, photo := range PhotosProject {
			newPhoto := models.PhotoProject{
				ProjectID: project.ID,
				URL:       photo,
			}

			_, err := h.ProjectRepositories.CreatePhotoProject(newPhoto)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "error 1"})
			}
		}
	}

	order, err := h.ProjectRepositories.GetOrderProjectByID(orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	order.Status = "Project is Complete"
	_, errOrder := h.ProjectRepositories.UpdateOrderProject(order)
	if errOrder != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errOrder.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataProject{
			Project: project,
		},
	})

}

func (h *handlerProject) GetProjectByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	project, err := h.ProjectRepositories.GetProjectByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataProject{
			Project: convertResponseProject(project),
		},
	})

}

func (h *handlerProject) GetProjectByOrderID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	project, err := h.ProjectRepositories.GetProjectByOrderID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataProject{
			Project: convertResponseProject(project),
		},
	})
}

func (h *handlerProject) UpdateProject(c echo.Context) error {
	description := c.FormValue("description")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	project, err := h.ProjectRepositories.GetProjectByOrderID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}
	if description != "" {
		project.Description = description
	}

	data, err := h.ProjectRepositories.UpdateProject(project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	photos, err := h.ProjectRepositories.GetPhotoProjectByProjectID(data.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}
	for _, photo := range photos {
		_, errDelete := h.ProjectRepositories.DeletePhoto(photo)
		if errDelete != nil {
			return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errDelete.Error()})
		}
	}

	var PhotosProject []string

	for i := 1; i <= 5; i++ {
		id := strconv.Itoa(i)
		image, err := c.FormFile("image" + id)
		if err == nil {
			src, err := image.Open()
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			defer src.Close()

			var ctx = context.Background()
			var CLOUD_NAME = os.Getenv("CLOUD_NAME")
			var API_KEY = os.Getenv("API_KEY")
			var API_SECRET = os.Getenv("API_SECRET")
			cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

			resp, errUpload := cld.Upload.Upload(ctx, src, uploader.UploadParams{Folder: "waysgallery"})
			if errUpload != nil {
				fmt.Println(errUpload.Error())
			}

			PhotosProject = append(PhotosProject, resp.SecureURL)
		}
	}

	if len(PhotosProject) != 0 {
		for _, photo := range PhotosProject {
			newPhoto := models.PhotoProject{
				ProjectID: project.ID,
				URL:       photo,
			}

			_, err := h.ProjectRepositories.CreatePhotoProject(newPhoto)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "error 1"})
			}
		}
	}
	order, err := h.ProjectRepositories.GetOrderProjectByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	order.Status = "Project is Complete"
	_, errOrder := h.ProjectRepositories.UpdateOrderProject(order)
	if errOrder != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: errOrder.Error()})
	}

	projectData, err := h.ProjectRepositories.GetProjectByID(data.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataProject{
			Project: convertResponseProject(projectData),
		},
	})

}

func convertResponseProject(p models.Project) projectdto.ProjectResponseDTO {
	return projectdto.ProjectResponseDTO{
		ID:          p.ID,
		Description: p.Description,
		Photos:      p.Photos,
		OrderID:     p.OrderID,
	}
}
