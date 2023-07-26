package handlers

import (
	"net/http"
	"strconv"
	bankdto "waysgallery/dto/bank"
	resultdto "waysgallery/dto/result"
	"waysgallery/models"
	"waysgallery/repositories"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type handlerBank struct {
	BankRepositories repositories.BankRepositories
}
type dataBanks struct {
	Banks interface{} `json:"banks"`
}

func HandlerBank(BankRepositories repositories.BankRepositories) *handlerBank {
	return &handlerBank{BankRepositories}
}

func (h *handlerBank) FindBanks(c echo.Context) error {
	banks, err := h.BankRepositories.FindBanks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataBanks{
			Banks: banks,
		},
	})
}

func (h *handlerBank) GetBankByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	bank, err := h.BankRepositories.GetBankByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataBanks{
			Banks: convertResponseBank(bank),
		},
	})
}

func (h *handlerBank) CreateBank(c echo.Context) error {
	request := new(bankdto.CreateBankDTO)

	errBind := c.Bind(&request)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: errBind.Error()})
	}

	validation := validator.New()

	validationErr := validation.Struct(request)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: validationErr.Error()})
	}

	newBank := models.Bank{
		Name: request.Name,
	}

	bank, err := h.BankRepositories.CreateBank(newBank)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: err.Error()})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataBanks{
			Banks: convertResponseBank(bank),
		},
	})
}

func (h *handlerBank) DeleteBank(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	bank, err := h.BankRepositories.GetBankByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	data, err := h.BankRepositories.DeleteBank(bank)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resultdto.ErrorResult{Status: "Failed", Message: "Invalid ID! Please input id as number."})
	}

	return c.JSON(http.StatusOK, resultdto.SuccessResult{
		Status: "Success",
		Data: dataBanks{
			Banks: convertResponseBank(data),
		},
	})
}

func convertResponseBank(b models.Bank) bankdto.BankResponseDTO {
	return bankdto.BankResponseDTO{
		ID:   b.ID,
		Name: b.Name,
	}
}
