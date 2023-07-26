package projectdto

import "waysgallery/models"

type ProjectResponseDTO struct {
	ID          int                   `json:"id"`
	Description string                `json:"description"`
	Photos      []models.PhotoProject `json:"photos"`
	OrderID     int                   `json:"order_id"`
}
